package nodes_test

import (
	builder "htestp/builder"
	httphandler "htestp/http_handler"
	"htestp/models"
	nodes "htestp/nodes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestDynamicNode struct {
	*nodes.DynamicNode
	*http.Client
}

func TestExecuteDynamicnode(t *testing.T) {

}
func TestGetNextNodes(t *testing.T) {
	t.Run("Added multiple nodes should be returned exactly when using GetNextNodes", func(t *testing.T) {
		StaticNodeOne := nodes.StaticNode{
			Request: httphandler.Request{
				Url:    "google.com",
				Method: string(models.GET),
			},
		}
		StaticNodeTwo := nodes.StaticNode{
			Request: httphandler.Request{
				Url:    "youtube.com",
				Method: string(models.GET),
			},
		}

		ExpectedNodes := []models.Node{&StaticNodeOne, &StaticNodeTwo}

		testBuilder := builder.CreateNewBuilder()
		testBuilder.AddDynamicNodeId("1", "trendyol.com", models.GET, nil, nil)
		testBuilder.AddStaticNodeBranch("google.com", models.GET, nil)
		testBuilder.AddStaticNodeBranch("youtube.com", models.GET, nil)

		assert.Equal(t, ExpectedNodes, (*testBuilder.Nodes)["1"].GetNextNodes())
	})
	// TODO: dont add nodes and check
}
func TestAddNode(t *testing.T) {
	t.Run("added dynamic node to the end by AddNode should match last node", func(t *testing.T) {
		testBuilder := builder.CreateNewBuilder()
		testBuilder.AddDynamicNodeId("1", "", models.GET, nil, nil)
		testNode := nodes.DynamicNode{
			InnerNode: nodes.StaticNode{
				Request: httphandler.Request{
					Url:    "https://dummyjson.com/products",
					Method: string(models.GET),
				},
			},
		}
		(*testBuilder.Nodes)["1"].AddNode(&testNode)

		assert.Equal(t,
			(*testBuilder.Nodes)["1"].(*nodes.DynamicNode).Next[0], &testNode)
	})
	t.Run("added nil to the end by AddNode should match nil", func(t *testing.T) {
		testBuilder := builder.CreateNewBuilder()
		testBuilder.AddDynamicNodeId("1", "", models.GET, nil, nil)

		(*testBuilder.Nodes)["1"].AddNode(nil)

		assert.Equal(t, (*testBuilder.Nodes)["1"].(*nodes.DynamicNode).Next[0], nil)
	})
}
func TestSuccessful(t *testing.T) {
	// TODO: Refactor to use testify
	t.Run("an unexecuted dynamic node should make the successful function return true", func(t *testing.T) {
		unExecutedBuilder := builder.CreateNewBuilder()
		unExecutedBuilder.AddDynamicNodeId("1", "https://dummyjson.com/products", models.GET, nil, nil).
			AddMatchConstraint("products[0].id", 1, models.TypeFloat)
		if !(*unExecutedBuilder.Nodes)["1"].Successful() {
			t.Errorf("Expected %t got %t", true, false)
		}
	})
	t.Run("an unsuccesful constraint should make the successful function return false", func(t *testing.T) {
		unsuccessfulBuilder := builder.CreateNewBuilder()
		unsuccessfulBuilder.AddDynamicNodeId("1", "https://dummyjson.com/products", models.GET, nil, nil).
			AddMatchConstraint("doesntexist.field", "welp", models.TypeString)
		unsuccessfulBuilder.Run()
		if (*unsuccessfulBuilder.Nodes)["1"].Successful() {
			t.Errorf("Expected %t got %t", false, true)
		}
	})
	t.Run("a successful constraint should make succesful function return true", func(t *testing.T) {
		successfulBuilder := builder.CreateNewBuilder()
		successfulBuilder.AddDynamicNodeId("1", "https://dummyjson.com/products", models.GET, nil, nil).
			AddMatchConstraint("products[0].id", 1, models.TypeFloat)
		successfulBuilder.Run()
		if !(*successfulBuilder.Nodes)["1"].Successful() {
			t.Errorf("Expected %t got %t", true, false)
		}
	})
}
