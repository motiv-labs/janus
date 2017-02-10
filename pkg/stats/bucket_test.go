package stats

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatsMetricName(t *testing.T) {
	dataProvider := []struct {
		Method string
		Path   string
		Metric string
	}{
		{"GET", "/", "get.-.-"},
		{"TRACE", "/api", "trace.api.-"},
		{"TRACE", "/api/", "trace.api.-"},
		{"POST", "/api/recipes", "post.api.recipes"},
		{"POST", "/api/recipes/", "post.api.recipes"},
		{"DELETE", "/api/recipes/123", "delete.api.recipes"},
		{"DELETE", "/api/recipes.foo-bar/123", "delete.api.recipes_foo-bar"},
		{"DELETE", "/api/recipes.foo_bar/123", "delete.api.recipes_foo__bar"},
		// paths withs IDs at the path second level
		{"GET", "/user/qwerty", "get.user." + pathIDPlaceholder},
		{"GET", "/users/qwerty", "get.users." + pathIDPlaceholder},
		{"GET", "/allergens/foobarbaz", "get.allergens." + pathIDPlaceholder},
		{"GET", "/cuisines/foobarbaz", "get.cuisines." + pathIDPlaceholder},
		{"GET", "/favorites/foobarbaz", "get.favorites." + pathIDPlaceholder},
		{"GET", "/ingredients/foobarbaz", "get.ingredients." + pathIDPlaceholder},
		{"GET", "/menus/foobarbaz", "get.menus." + pathIDPlaceholder},
		{"GET", "/ratings/foobarbaz", "get.ratings." + pathIDPlaceholder},
		{"GET", "/recipes/foobarbaz", "get.recipes." + pathIDPlaceholder},
		{"GET", "/addresses/foobarbaz", "get.addresses." + pathIDPlaceholder},
		{"GET", "/boxes/foobarbaz", "get.boxes." + pathIDPlaceholder},
		{"GET", "/coupons/foobarbaz", "get.coupons." + pathIDPlaceholder},
		{"GET", "/customers/foobarbaz", "get.customers." + pathIDPlaceholder},
		{"GET", "/delivery_options/foobarbaz", "get.delivery__options." + pathIDPlaceholder},
		{"GET", "/product_families/foobarbaz", "get.product__families." + pathIDPlaceholder},
		{"GET", "/products/foobarbaz", "get.products." + pathIDPlaceholder},
		{"GET", "/recipients/foobarbaz", "get.recipients." + pathIDPlaceholder},
		// path may have either numeric ID or non-numeric trackable path
		{"GET", "/subscriptions/12345", "get.subscriptions." + pathIDPlaceholder},
		{"GET", "/subscriptions/search", "get.subscriptions.search"},
		{"GET", "/freebies/12345", "get.freebies." + pathIDPlaceholder},
		{"GET", "/freebies/search", "get.freebies.search"},
	}

	for _, data := range dataProvider {
		assert.Equal(t, data.Metric, getMetricName(data.Method, &url.URL{Path: data.Path}))
	}
}
