# konfa

## TODO:

Setup CLI Tool for interacting with application.
Add following functionality:

- Update configmap by specifying key and new value, instead of having to use Vim for jumping in and editing by hand.
- Update deployment related keys with new values, e.g `konfa deployment replica newVal`, where newVal is an integer.
- In general updates different kubernetes related config in a simpler way than having to use kubectl and specifying.
- Also add possibility of setting different contexts for the application.
- `konfa set namespace <name>`, and therefor don't have to specify namespace anymore when updating certain kubernetes config.
- Add functionality for some kind of diff viewer, where the newly added config will show up with some nice color coding to show that entry X -> Y.
