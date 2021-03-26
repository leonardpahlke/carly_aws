package shared

import (
	"encoding/json"

	pkg "github.com/leonardpahlke/carly_pkg"

	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/iam"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

const none = pulumi.String("")
const IamPolicy_Action_AssumeRole = pulumi.String("sts:AssumeRole")
const IamPolicy_Effect_Allow = pulumi.String("Allow")
const IamPolicy_Resource_All = pulumi.String("*")
const IamPolicy_Resource_None = none
const IamPolicy_Principal_None = none
const IamPolicy_Sid_None = none
const IamPolicy_Service_Lambda = pulumi.String("lambda.amazonaws.com")

// create a policy statement with preset variables - Actions can be set
func IamCreatePolicyString_Actions(actions pulumi.StringArray) pulumi.String {
	policyStatement := ownIamPolicyStatement{}
	policyStatement.addActions(actions)
	policyStatement.addEffect(IamPolicy_Effect_Allow)
	policyStatement.addResource(IamPolicy_Resource_All)
	return policyStatement.getPolicyStatement()
}

// create a policy statement with preset variables - Actions and Resources can be set
func IamCreatePolicyString_Actions_Resource(actions pulumi.StringArray, resource pulumi.String) pulumi.String {
	policyStatement := ownIamPolicyStatement{}
	policyStatement.addActions(actions)
	policyStatement.addEffect(IamPolicy_Effect_Allow)
	policyStatement.addResource(resource)
	return policyStatement.getPolicyStatement()
}

func IamCreatePolicyString_Actions_Resource_StrOut(actions pulumi.StringArray, resource pulumi.StringOutput) pulumi.String {
	policyStatement := ownIamPolicyStatement{}
	policyStatement.addActions(actions)
	policyStatement.addEffect(IamPolicy_Effect_Allow)
	policyStatement.addResourceStrOut(resource)
	return policyStatement.getPolicyStatement()
}

// create a policy statement with preset variables - Principal can be set
func IamCreatePolicyString_Assume_Policy(principalService pulumi.String) pulumi.StringOutput {
	return pulumi.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Principal": {
					"Service": "%s"
				},
				"Action": "sts:AssumeRole"
			}
		]
	}`, principalService)
	// policyStatement := ownIamPolicyStatement{}
	// policyStatement.addAction(IamPolicy_Action_AssumeRole)
	// policyStatement.addEffect(IamPolicy_Effect_Allow)
	// policyStatement.addPrincipal(principalService)
	// return policyStatement.getPolicyStatement()
}

type ownIamPolicyStatement struct {
	policy pulumi.Map
}

func (ps ownIamPolicyStatement) addActions(actions pulumi.StringArray) ownIamPolicyStatement {
	ps.policy["Action"] = actions
	return ps
}

func (ps ownIamPolicyStatement) addAction(action pulumi.String) ownIamPolicyStatement {
	ps.policy["Action"] = action
	return ps
}

func (ps ownIamPolicyStatement) addEffect(effect pulumi.String) ownIamPolicyStatement {
	ps.policy["Effect"] = effect
	return ps
}

func (ps ownIamPolicyStatement) addResource(resource pulumi.String) ownIamPolicyStatement {
	ps.policy["Resource"] = resource
	return ps
}

func (ps ownIamPolicyStatement) addResourceStrOut(resource pulumi.StringOutput) ownIamPolicyStatement {
	ps.policy["Resource"] = resource
	return ps
}

func (ps ownIamPolicyStatement) addSid(sid pulumi.String) ownIamPolicyStatement {
	ps.policy["Sid"] = sid
	return ps
}

func (ps ownIamPolicyStatement) addPrincipal(servicePrincipal pulumi.String) ownIamPolicyStatement {
	tmpJSON, err := json.Marshal(map[string]interface{}{
		"Service": servicePrincipal,
	})
	if err != nil {
		pkg.LogError("OwnIamPolicyStatement.addPrincipal", "could not marshal policy principal statement object", err)
	}
	ps.policy["Principal"] = pulumi.String(string(tmpJSON))
	return ps
}

func (ps ownIamPolicyStatement) getPolicyStatement() pulumi.String {
	statement := pulumi.MapArray{ps.policy}

	// create JSON object
	tmpJSON, err := json.Marshal(pulumi.Map{
		"Version":   pulumi.String("2012-10-17"),
		"Statement": statement,
	})
	if err != nil {
		pkg.LogError("OwnIamPolicyStatement.getPolicyStatement", "could not marshal policy statement object", err)
	}
	return pulumi.String(string(tmpJSON))
}

/* create a policy statement
func IamCreatePolicyString(actions pulumi.StringArray, singleAction pulumi.String, effect pulumi.String, resource pulumi.String, principal pulumi.String, sid pulumi.String, omitResource bool) pulumi.String {
	// create an empty statement
	tempStatement := pulumi.Map{}

	// add fields to statement if information is set
	if len(actions) != 0 {
		tempStatement["Action"] = actions
	} else if singleAction != none {
		tempStatement["Action"] = singleAction
	}
	if effect != none {
		tempStatement["Effect"] = effect
	}
	if principal != none {
		tempStatement["Principal"] = principal
	}
	if pulumi.String(sid.ElementType().String()) != none {
		tempStatement["Sid"] = sid
	}
	if !omitResource {
		tempStatement["Resource"] = resource
	}

	statement := pulumi.MapArray{tempStatement}

	// create JSON object
	tmpJSON, err := json.Marshal(pulumi.Map{
		"Version":   pulumi.String("2012-10-17"),
		"Statement": statement,
	})
	if err != nil {
		LogError("pkg.CreatePolicyString", "could not marshal policy statement object", err)
	}
	return pulumi.String(string(tmpJSON))
}
*/

/*
	POLICY
*/

// Create policy

// create inline-policy
func CreateInlinePolicyStatement(name string, policyString pulumi.StringOutput) iam.RoleInlinePolicyInput {
	return iam.RoleInlinePolicyArgs{
		Name:   pulumi.Sprintf("inline-policy-%s", name),
		Policy: policyString,
	}
}

/*
	ROLE
*/
func CreateRole(ctx *pulumi.Context, roleName string, policyAssumeStringJson pulumi.StringOutput, inlinePolicyAry iam.RoleInlinePolicyArray) *iam.Role {
	role, err := iam.NewRole(
		ctx,
		GetResourceName(roleName),
		&iam.RoleArgs{
			AssumeRolePolicy: policyAssumeStringJson,
			InlinePolicies:   inlinePolicyAry,
			Tags:             GetTags(roleName),
		})
	if err != nil {
		pkg.LogError("pkg.CreateRole", "could not create role", err)
	}
	return role
}
