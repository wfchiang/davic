# DAVIC: Structured Data as Program

**Davic** framework allows you to transit **computations** in networks. Cool? (Nah... I think it is cool.)

Here is the idea: write programs in JSON format (a structured data format). So, your programs are data now. The truth is, programs are data, period. 

The bright side of writing/encoding programs in JSON format is that we can ride all the benefits of the RESTful service model. 
Now one can easily transit her/his programs among internet nodes (or, in cloud). It can mean many things.
But, obviously, program + data basically defines a computation. 
So, you can now transit your computation to somewhere else (like coud). 

Here are some example apps which demonstrate how Davic works and what it can achieve. 
* [REST-mocker](https://github.com/wfchiang/davic/tree/master/src/sample-apps/rest-mocker) A mock RESTful service which can be defined "what to mock" **at runtime**. Yes, you heart me. **At runtime!**