# Contributing

Contributions are **welcome** and will be fully **credited**.

We accept contributions via Pull Requests on [Github](https://github.com/hellofresh/location-service).

## Pull Requests

- **Add tests!** - Your patch won't be accepted if it doesn't have tests.

- **Document any change in behaviour** - Make sure the `README.md` and any other relevant documentation are kept up-to-date.

- **Create feature branches** - Don't ask us to pull from your master branch.

- **One pull request per feature** - If you want to do more than one thing, send multiple pull requests.

- **Send coherent history** - Make sure each individual commit in your pull request is meaningful. If you had to make multiple intermediate commits while developing, please squash them before submitting.

## Deployments

Deployments on the location service follows the current flow:

- Any feature that is merged into master a CI build will be triggered and deployed to **staging** automatically (If the tests pass)

- If the Product owner and Team Lead agreed that what we have on master is enough to go live, we create a release (tag)

- After a tag is created a CI build will be triggered and deployed to **live** automatically (If the tests pass)

- **Consider our release cycle** - We try to follow [SemVer v2.0.0](http://semver.org/). Randomly breaking public APIs is not an option.

## Running Tests

``` bash
$ go test -v ./...
```

**Happy coding**!
