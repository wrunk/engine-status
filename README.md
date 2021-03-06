# Engine Status

an App Engine server running and you want to gain more insight
into how many unique instances and versions are running then
this is the tool for you! Shows you:

- How many different versions of your service are running
- How many total instances are running
- Which data centers are these served from

This tool is helpful to run before and during a deploy.

**Warning this tool can rack up tremendous \$$$charges$$\$ if used
improperly. Please use with caution.** It cannot be used to ddos
a random endpoint since it requires a proper return format.

## Quick start

1. Setup a server endpoint as described below
1. `go run *.go <url> <key>`

## Server Code

In order for this to work properly, you need to create a JSON endpoint
to be deployed with your app engine app. Here is an example using echo
webserver:

```go
package main

import (
    os
)

// An endpoint for use by the engine-status tool
func engineStatus(ctx echo.Context) error {
	sleepTimeMS := ctx.QueryParam("s")
	st, err := strconv.Atoi(sleepTimeMS)
	if err == nil {
		fmt.Printf("Sleeping for (%d) MS\n", st)
		time.Sleep(time.Duration(st))
	}
	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"gae_application":   os.Getenv("GAE_APPLICATION"),
		"gae_deployment_id": os.Getenv("GAE_DEPLOYMENT_ID"),
		"gae_env":           os.Getenv("GAE_ENV"),
		"gae_instance":      os.Getenv("GAE_INSTANCE"),
		"gae_memory_mb":     os.Getenv("GAE_MEMORY_MB"),
		"gae_service":       os.Getenv("GAE_SERVICE"),
		"gae_runtime":       os.Getenv("GAE_RUNTIME"),
		"gae_version":       os.Getenv("GAE_VERSION"),

		"google_cloud_project": os.Getenv("GOOGLE_CLOUD_PROJECT"),
		// Only used with new, second generation runtimes
		"port": os.Getenv("PORT"),
	})
}
...

```

## Similar To

- Apache Bench (ab)
- App Engine dashboard(s)
- Load testing tools
