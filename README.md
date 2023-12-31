# konfa

Go application for interacting with k8s cluster, without having to use kubectl edit.

#### **Heads up:**

This is still a work in progress and in early experimental stage, don't try this on any live cluster. Been doing most of my tests on a local k3d cluster.

### Verify your changes:

For whole deployment:

```
kubectl get deployment -n <namespace> -o yaml
```

For just configmap:

```
kubectl get configmaps -n <namespace> <configmapname> -o yaml
```

## TODO:

Setup CLI Tool for interacting with application.
Add following functionality:

- [High prio] Make this to a command so people can simply run konfa instead of running the source code and go run main.go
- Update configmap by specifying key and new value, instead of having to use Vim for jumping in and editing by hand. [Done]
- Update deployment related keys with new values, e.g `konfa deployment replica newVal`, where newVal is an integer. [Done]
- In general updates different kubernetes related config in a simpler way than having to use kubectl and specifying. [Done]
- `konfa set namespace <name>`, and therefore don't have to specify namespace anymore when updating certain kubernetes config. [Done]
- Add functionality for some kind of diff viewer, where the newly added config will show up with some nice color coding to show that entry X -> Y.
- Add some functionality for fetching config in a simpler way than having to run kubectl commands?
- Add tests ofc

## Usage

Easy tool for running incident rehearsal / preperation, for quickly re-configuring cluster data without having to use kubectl edit + rollout restart.

### Testing out locally:

Setup cluster:

```
k3d cluster create <yourclustername>
```

Example commands:

```
go run main.go config set <key> <value> --is-file=true --config=myconfig
```

```
go run main.go container nginx set resources.requests 400Mi
```

```
go run main.go --help
```

## Future additions:

- Make it look nicer, maybe look into some lib like bubbletea?
