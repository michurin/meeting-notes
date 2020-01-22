```
go test -v .
```

**Appendix A:** I used here `github.com/looplab/fsm` only for historical reasons.
I find it slightly weird.

- Inconvenient way to write STF
- Dangerous methods like SetState. You have direct access to FSM internals
- Callbacks triggers in weird way
- The engine knows nothing about final states, you can not specify them explicitly and have to treat them manually in your code
- Fancy naming. The same thing can be have two names like `Name` and `Event`