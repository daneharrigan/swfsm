package migrator

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/gen/kinesis"
	"github.com/awslabs/aws-sdk-go/gen/swf"
)

var testDomain = fmt.Sprintf("test-domain-%d", time.Now().UnixNano())

func TestMigrateDomains(t *testing.T) {
	if os.Getenv("AWS_ACCESS_KEY_ID") == "" || os.Getenv("AWS_SECRET_ACCESS_KEY") == "" {
		log.Printf("WARNING: NO AWS CREDS SPECIFIED, SKIPPING MIGRATIONS TEST")
		return
	}

	creds, _ := aws.EnvCreds()
	client := swf.New(creds, "us-east-1", nil)

	domain := fmt.Sprintf("test-domain-%d", time.Now().UnixNano()) //dont use the testDomain since we deprecate this one

	req := swf.RegisterDomainInput{
		Name:                                   aws.String(domain),
		Description:                            aws.String("test domain"),
		WorkflowExecutionRetentionPeriodInDays: aws.String("30"),
	}

	d := DomainMigrator{
		RegisteredDomains: []swf.RegisterDomainInput{req},
		Client:            client,
	}

	d.Migrate()
	d.Migrate()

	dep := swf.DeprecateDomainInput{
		Name: aws.String(domain),
	}

	dd := DomainMigrator{
		DeprecatedDomains: []swf.DeprecateDomainInput{dep},
		Client:            client,
	}

	dd.Migrate()
	dd.Migrate()

}

func TestMigrateWorkflowTypes(t *testing.T) {
	if os.Getenv("AWS_ACCESS_KEY_ID") == "" || os.Getenv("AWS_SECRET_ACCESS_KEY") == "" {
		log.Printf("WARNING: NO AWS CREDS SPECIFIED, SKIPPING MIGRATIONS TEST")
		return
	}
	createDomain()
	creds, _ := aws.EnvCreds()
	client := swf.New(creds, "us-east-1", nil)

	workflow := fmt.Sprintf("test-workflow-%d", time.Now().UnixNano())
	version := fmt.Sprintf("test-workflow-version-%d", time.Now().UnixNano())

	req := swf.RegisterWorkflowTypeInput{
		Name:        &workflow,
		Description: aws.String("test workflow migration"),
		Version:     &version,
		Domain:      aws.String(testDomain),
	}

	w := WorkflowTypeMigrator{
		RegisteredWorkflowTypes: []swf.RegisterWorkflowTypeInput{req},
		Client:                  client,
	}

	w.Migrate()
	w.Migrate()

	dep := swf.DeprecateWorkflowTypeInput{
		WorkflowType: &swf.WorkflowType{
			Name:    aws.String(workflow),
			Version: aws.String(version),
		},
		Domain: aws.String(testDomain),
	}

	wd := WorkflowTypeMigrator{
		DeprecatedWorkflowTypes: []swf.DeprecateWorkflowTypeInput{dep},
		Client:                  client,
	}

	wd.Migrate()
	wd.Migrate()

}

func TestMigrateActivityTypes(t *testing.T) {

	if os.Getenv("AWS_ACCESS_KEY_ID") == "" || os.Getenv("AWS_SECRET_ACCESS_KEY") == "" {
		log.Printf("WARNING: NO AWS CREDS SPECIFIED, SKIPPING MIGRATIONS TEST")
		return
	}
	createDomain()
	creds, _ := aws.EnvCreds()
	client := swf.New(creds, "us-east-1", nil)

	activity := fmt.Sprintf("test-activity-%d", time.Now().UnixNano())
	version := fmt.Sprintf("test-activity-version-%d", time.Now().UnixNano())

	req := swf.RegisterActivityTypeInput{
		Name:        &activity,
		Description: aws.String("test activity migration"),
		Version:     &version,
		Domain:      aws.String(testDomain),
	}

	a := ActivityTypeMigrator{
		RegisteredActivityTypes: []swf.RegisterActivityTypeInput{req},
		Client:                  client,
	}

	a.Migrate()
	a.Migrate()

	dep := swf.DeprecateActivityTypeInput{
		ActivityType: &swf.ActivityType{
			Name:    &activity,
			Version: &version,
		},
		Domain: aws.String(testDomain),
	}

	ad := ActivityTypeMigrator{
		DeprecatedActivityTypes: []swf.DeprecateActivityTypeInput{dep},
		Client:                  client,
	}

	ad.Migrate()
	ad.Migrate()

}

func TestMigrateStreams(t *testing.T) {
	if os.Getenv("AWS_ACCESS_KEY_ID") == "" || os.Getenv("AWS_SECRET_ACCESS_KEY") == "" {
		log.Printf("WARNING: NO AWS CREDS SPECIFIED, SKIPPING MIGRATIONS TEST")
		return
	}

	creds, _ := aws.EnvCreds()
	client := kinesis.New(creds, "us-east-1", nil)

	sm := StreamMigrator{
		Streams: []kinesis.CreateStreamInput{
			kinesis.CreateStreamInput{
				StreamName: aws.String(testDomain),
				ShardCount: aws.Integer(1),
			},
		},
		Client: client,
	}

	sm.Migrate()
	sm.Migrate()

}

func createDomain() {
	creds, _ := aws.EnvCreds()
	client := swf.New(creds, "us-east-1", nil)
	req := swf.RegisterDomainInput{
		Name:                                   aws.String(testDomain),
		Description:                            aws.String("test domain"),
		WorkflowExecutionRetentionPeriodInDays: aws.String("30"),
	}
	w := DomainMigrator{
		RegisteredDomains: []swf.RegisterDomainInput{req},
		Client:            client,
	}

	w.Migrate()
}
