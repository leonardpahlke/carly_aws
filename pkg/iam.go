package pkg

//const IAMPolicyActionsAssumeRole = "sts:AssumeRole"
//const IAMPolicyActionsS3PutObject = "s3:PutObject"
//const IAMPolicyPrincipalsIdentifiersEc2 = "ec2.amazonaws.com"
//const IAMPolicyPrincipalsIdentifiersLambda = "lambda.amazonaws.com"
//const IAMPolicyEffectAllow = "Allow"
//
///*
//AssumeRolePolicy: pulumi.String(`{
//				"Version": "2012-10-17",
//				"Statement": [{
//					"Sid": "",
//					"Effect": "Allow",
//					"Principal": {
//						"Service": "lambda.amazonaws.com"
//					},
//					"Action": "sts:AssumeRole"
//				}]
//			}`),
// */
//
//func CreateRole(ctx *pulumi.Context) (*iam.Role, error) {
//	instance_assume_role_policy, err := iam.GetPolicyDocument(ctx, &iam.GetPolicyDocumentArgs{
//		Statements: []iam.GetPolicyDocumentStatement{
//			iam.GetPolicyDocumentStatement{
//				Actions: []string{
//					"sts:AssumeRole",
//				},
//				Principals: []iam.GetPolicyDocumentStatementPrincipal{
//					iam.GetPolicyDocumentStatementPrincipal{
//						Type: "Service",
//						Identifiers: []string{
//							"ec2.amazonaws.com",
//						},
//					},
//				},
//			},
//		},
//	}, nil)
//	if err != nil {
//		return &iam.Role{}, err
//	}
//	createdRole, err := iam.NewRole(ctx, "instance", &iam.RoleArgs{
//		Path:             pulumi.String("/system/"),
//		AssumeRolePolicy: pulumi.String(instance_assume_role_policy.Json),
//	})
//	return createdRole, nil
//}
//
//func CreatePolicyDocument(actions []string, identifiers []string) iam.GetPolicyDocumentStatement {
//	return iam.GetPolicyDocumentStatement{
//		Actions: actions,
//		Principals: []iam.GetPolicyDocumentStatementPrincipal{
//			iam.GetPolicyDocumentStatementPrincipal{
//				Type: "Service",
//				Identifiers: identifiers,
//			},
//		},
//	}
//}
//
//func CreatePrinciplePolicyDocument(identifiers []string) iam.GetPolicyDocumentStatementPrincipal {
//	return iam.GetPolicyDocumentStatementPrincipal{
//		Type: "Service",
//		Identifiers: identifiers,
//		// identifier example: "ec2.amazonaws.com"
//	}
//}
//
//type IamPolicyStatement struct {
//	Sid string
//	Effect string
//	Principal []IamPrincipal
//}
//
//
//type IamPrincipal struct {
//	Service string
//}
//
//
//func CreatePolicyStatement(policyStatement []IamPolicyStatement) string {
//	return CreateKeyValuePairs(structToMap(policyStatement))
//}
//
//func createSinglePolicyStatement(policyStatement IamPolicyStatement) string {
//	return ""
//}
//
///*
//This function will help you to convert your object from struct to map[string]interface{} based on your JSON tag in your structs.
//Example how to use posted in sample_test.go file.
//*/
//func structToMap(item interface{}) map[string]interface{} {
//
//	res := map[string]interface{}{}
//	if item == nil {
//		return res
//	}
//	v := reflect.TypeOf(item)
//	reflectValue := reflect.ValueOf(item)
//	reflectValue = reflect.Indirect(reflectValue)
//
//	if v.Kind() == reflect.Ptr {
//		v = v.Elem()
//	}
//	for i := 0; i < v.NumField(); i++ {
//		tag := v.Field(i).Tag.Get("json")
//		field := reflectValue.Field(i).Interface()
//		if tag != "" && tag != "-" {
//			if v.Field(i).Type.Kind() == reflect.Struct {
//				res[tag] = structToMap(field)
//			} else {
//				res[tag] = field
//			}
//		}
//	}
//	return res
//}
//
//// Effect
//const ConstIamEffectAllow = "Allow"
//
//// Principle
//const ConstIamPrincipleServiceLambda = "lambda.amazonaws.com"
//
//// Action
//const ConstIamActionAssumeRole = "sts:AssumeRole"
//
//
//
///*
//pulumi.String(`{
//				"Version": "2012-10-17",
//				"Statement": [{
//					"Sid": "",
//					"Effect": "Allow",
//					"Principal": {
//						"Service": "lambda.amazonaws.com"
//					},
//					"Action": "sts:AssumeRole"
//				}]
//			}`
// */