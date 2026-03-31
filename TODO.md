# TODO

Items to create as GitHub issues when the repository is published.

## Parser Limitations

- [ ] $(MAKE) detection: target is a variable (`$(MAKE) $(GOALS)`) — no variable expansion attempted
- [ ] $(MAKE) detection: complex commands with pipes or redirects (`$(MAKE) target | grep error`)
- [ ] $(MAKE) detection: loops (`for t in $(TARGETS); do $(MAKE) $$t; done`)
- [ ] $(MAKE) detection: multiple make invocations on one line (only the first is detected)
