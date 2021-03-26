package internal

import (
	"bytes"
	"carly_aws/shared"
	"fmt"
	"strings"

	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/sfn"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

const methodDelete = "delete"
const methodGet = "get"
const methodPut = "put"
const methodUpdate = "update"

func CreateSpiderStateMachine(ctx *pulumi.Context, config SpiderStateMachineConfig) (SpiderStateMachineData, error) {
	_, err := sfn.NewStateMachine(ctx, "SpiderStateMachine", &sfn.StateMachineArgs{
		Name:       pulumi.Sprintf("SpiderStateMachine"),
		RoleArn:    pulumi.String(""), // todo step function role arn
		Definition: pulumi.String(fmt.Sprintf("%v%v%v%v%v%v%v%v%v%v%v%v%v", "{\n", "  \"Comment\": \"A Hello World example of the Amazon States Language using an AWS Lambda Function\",\n", "  \"StartAt\": \"HelloWorld\",\n", "  \"States\": {\n", "    \"HelloWorld\": {\n", "      \"Type\": \"Task\",\n", "      \"Resource\": \"", "aws_lambda_function.Lambda.Arn", "\",\n", "      \"End\": true\n", "    }\n", "  }\n", "}\n")),
		Tags:       shared.GetTags("SpiderStateMachine"),
	})
	if err != nil {
		return SpiderStateMachineData{}, err
	}
	return SpiderStateMachineData{}, nil
}

type SpiderStateMachineConfig struct{}

type SpiderStateMachineData struct{}

func createKeyValuePairs(m map[string]string) string {
	b := new(bytes.Buffer)
	valEscapeString := ""
	_, _ = fmt.Fprint(b, "{\n  ")
	for key, value := range m {
		if strings.ContainsAny(value, "{}") {
			valEscapeString = ""
		} else {
			valEscapeString = "\""
		}
		_, _ = fmt.Fprintf(b, "  \"%s\": %s%s%s\n  ", key, valEscapeString, value, valEscapeString)
	}
	_, _ = fmt.Fprint(b, "}")
	return b.String()
}

// Tasks definitions
/*
"Invoke Lambda function": {
  "Type": "Task",
  "Resource": "arn:aws:states:::lambda:invoke",
  "Parameters": {
    "FunctionName": "arn:aws:lambda:REGION:ACCOUNT_ID:function:FUNCTION_NAME",
    "Payload": {
      "Input.$": "$"
    }
  },
  "Next": "NEXT_STATE"
}
*/
func createTaskLambdaFunction(taskName string, lambdaArn string, functionName string, payload map[string]string, nextState string) string {
	lambdaFuncTaskDefinition := createKeyValuePairs(map[string]string{
		taskName: createKeyValuePairs(map[string]string{
			"Type":     "Task",
			"Resource": lambdaArn,
			"Next":     nextState,
			"Parameters": createKeyValuePairs(map[string]string{
				"FunctionName": functionName,
				"Payload":      createKeyValuePairs(payload),
			}),
		}),
	})
	return lambdaFuncTaskDefinition
}

/*
"Put item into DynamoDB": {
  "Type": "Task",
  "Resource": "arn:aws:states:::dynamodb:[putItem / deleteItem]",
  "Parameters": {
    "TableName": "MyDynamoDBTable",
    ["Item": {
      "Column": {
        "S": "MyEntry"
      }
    },
	"Key": {
      "Column": {
        "S": "MyEntry"
      }
    }
	]
  },
  "Next": "NEXT_STATE"
}
*/
func createTaskDynamoDB(method string, tableName string, updateInformation UpdateInformation, itemKeyValue map[string]string, nextState string) string {
	resource := ""
	parameterKey := "Key"
	switch method {
	case methodDelete:
		resource = "arn:aws:states:::dynamodb:deleteItem"
	case methodGet:
		resource = "arn:aws:states:::dynamodb:getItem"
	case methodPut:
		resource = "arn:aws:states:::dynamodb:putItem"
		parameterKey = "Item"
	case methodUpdate:
		resource = "arn:aws:states:::dynamodb:updateItem"
	}
	if resource == "" {
		// raise exception
	}

	dynamoDbTaskDefinition := ""
	if resource == methodUpdate {
		dynamoDbTaskDefinition = createKeyValuePairs(map[string]string{
			"Type":     "Task",
			"Resource": resource,
			"Next":     nextState,
			"Parameters": createKeyValuePairs(map[string]string{
				"UpdateExpression":          updateInformation.updateExpression,
				"ExpressionAttributeValues": createKeyValuePairs(updateInformation.expressionAttributeValues),
				"TableName":                 tableName,
				parameterKey:                createKeyValuePairs(itemKeyValue),
			}),
		})
	} else {
		dynamoDbTaskDefinition = createKeyValuePairs(map[string]string{
			"Type":     "Task",
			"Resource": resource,
			"Next":     nextState,
			"Parameters": createKeyValuePairs(map[string]string{
				"TableName":  tableName,
				parameterKey: createKeyValuePairs(itemKeyValue),
			}),
		})
	}

	return dynamoDbTaskDefinition
}

type UpdateInformation struct {
	updateExpression          string
	expressionAttributeValues map[string]string
}

/*
"Update item into DynamoDB": {
  "Type": "Task",
  "Resource": "arn:aws:states:::dynamodb:updateItem",
  "Parameters": {
    "TableName": "MyDynamoDBTable",
    "Key": {
      "Column": {
        "S": "MyEntry"
      }
    },
    "UpdateExpression": "SET MyKey = :myValueRef",
    "ExpressionAttributeValues": {
      ":myValueRef": {
        "S": "MyValue"
      }
    }
  },
  "Next": "NEXT_STATE"
}
*/

/*
"Choice State": {
  "Type": "Choice",
  "Choices": [
    {
      "Not": {
        "Variable": "$.type",
        "StringEquals": "Private"
      },
      "Next": "NEXT_STATE_ONE"
    },
    {
      "Variable": "$.value",
      "NumericEquals": 0,
      "Next": "NEXT_STATE_TWO"
    },
    {
      "And": [
        {
          "Variable": "$.value",
          "NumericGreaterThanEquals": 20
        },
        {
          "Variable": "$.value",
          "NumericLessThan": 30
        }
      ],
      "Next": "NEXT_STATE_THREE"
    }
  ],
  "Default": "DEFAULT_STATE"
}
*/
func getChoiceParserStepFunction(variableToCompare string, newspaperIdentifier []newsPaperParserInfo) string {
	choiceStr := "\"Type\": \"Choice\", " +
		"\"Choices\": [{ "
	for i, e := range newspaperIdentifier {
		choiceStr += strEqualsChoiceStepFunction(variableToCompare, e.id, e.nextState)
		if i != 0 {
			choiceStr += ","
		}
	}
	choiceStr += "}]"
	return choiceStr
	//return "\"Type\": \"Choice\", " +
	//	   "\"Choices\": [{ " +
	//			"\"Not\": { " +
	//				"\"Variable\": \"$.type\", " +
	//				"\"StringEquals\": \"Private\" " +
	//			"}, " +
	//			"\"Next\": \"NEXT_STATE_ONE\"}," +
	//
	//			"{\"Variable\": \"$.value\"," +
	//			"\"NumericEquals\": 0," +
	//			"\"Next\": \"NEXT_STATE_TWO\"}," +
	//
	//		  	"{\"And\": [" +
	//				"{\"Variable\": \"$.value\"," +
	//				"\"NumericGreaterThanEquals\": 20}, " +
	//				"{\"Variable\": \"$.value\"," +
	//				"\"NumericLessThan\": 30}" +
	//			"]," +
	//			"\"Next\": \"NEXT_STATE_THREE\"}" +
	//		"]," +
	//
	//		"\"Default\": \"DEFAULT_STATE\""
}

func strEqualsChoiceStepFunction(variable string, stringEquals string, nextState string) string {
	return "{" +
		"\"Variable\": \"$." + variable + "\"," +
		"\"StringEquals\": " + stringEquals +
		"\"Next\": \"" + nextState + "\"" +
		"}"
}

type newsPaperParserInfo struct {
	id        string
	nextState string
}
