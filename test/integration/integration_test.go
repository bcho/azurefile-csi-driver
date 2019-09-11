/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package integration

import (
	"context"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/kubernetes-sigs/azurefile-csi-driver/test/azure"
	"github.com/kubernetes-sigs/azurefile-csi-driver/test/credentials"
	"github.com/stretchr/testify/assert"
)

func TestIntegration(t *testing.T) {
	creds, err := credentials.Get()
	assert.NoError(t, err)
	assert.NotNil(t, creds)

	os.Setenv("AZURE_CREDENTIAL_FILE", credentials.TempAzureCredentialFilePath)

	azureClient, err := azure.GetAzureClient(creds.SubscriptionID, creds.AADClientID, creds.TenantID, creds.AADClientSecret)
	assert.NoError(t, err)

	ctx := context.Background()
	// Create an empty resource group for integration test
	t.Logf("Creating resource group %s", creds.ResourceGroup)
	_, err = azureClient.EnsureResourceGroup(ctx, creds.ResourceGroup, creds.Location, nil)
	assert.NoError(t, err)
	defer func() {
		t.Logf("Deleting resource group %s", creds.ResourceGroup)
		err := azureClient.DeleteResourceGroup(ctx, creds.ResourceGroup)
		assert.NoError(t, err)
	}()

	// Execute the script from project root
	err = os.Chdir("../..")
	assert.NoError(t, err)

	cwd, err := os.Getwd()
	assert.NoError(t, err)
	assert.True(t, strings.HasSuffix(cwd, "azurefile-csi-driver"))

	cmd := exec.Command("./test/integration/run-tests-all-clouds.sh")
	cmd.Dir = cwd
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("Integration test failed %v", err)
	}
}
