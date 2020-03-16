# QA

q: What was the estimation?

a: 4hrs. Why? Because I wasn't accurate when read. I wasted 15 hours to do everything there...

---

q: what is the structure of all repos?

a: api - separate repo for api, processor - stream processor repo, not separated repos because I don't want to remove multiple repos(3) from my account after review.

---

q: How to build service processor(first service that process stream of data)?

a: 
    
    make

---

q: How to run?

a: realisation not done(there is no second service). But you can just launch binary or docker.

---

q: How to run tests?

a: 

    make test-unit
    make test-integration

---

q: Why not done?

a: I'm really tired to write services for more than 15 hours just for test task.

---

q: What is done?

a:  
- First service core functionality
- small part of service unit tests
- one benchmark
- small part of service integration tests
- showed usage of concurrency, atomics, different sync methodoligies, stream parsing...
- DRY
- go.mod, Makefile(CI/CD), docker...

And a lot more...

---

q: Why config is static(not env)? 

a: it's just a test task

---

q: What is not done?

a: 
- Readiness probe
- Proper config load
- Metrics
- Logs
- Statistics
- proper jobs Queueing
- get endpoint
- second service(storage)
- circuit breaker
- a lot of tests
- linter - one line with config? Please no, I don't wan't to play with it now 
- docker-compose - one more thing to do, it's just a small file...

...


and this list also is not done) 

There is a lot of small pieces missed, because it will not be deployed anyway. So why not just ask questions related to missed things and discuss them?
