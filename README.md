# DAVIC: Structured Data as Program

In this cloud era, computations are done by microservices with data transitions among them. We typically follow these steps to build a system: 
1. Define computation units
2. Define data structures (e.g., JSON schemas) 
3. Implement the computation units as microservices 
4. Deploy the microservices to cloud 

Whenever you want to change the behaviors of the system, add new functionalities, or fix defects, you will need to re-implement the microservices, and re-deploy them to cloud. 
This is not flexible to adapt the fast-pace changing of the current needs for cloud. 

Here is the **Davic**-way to build a system: 
1. Define computation units 
2. Define data structures 
3. Deploy a set of **Davic** generic microservices 
4. Through network, send the computation tasks to the **Davic** at runtime. Yes. I am saying "sending computations as data (e.g., JSON objects)." 

Whenever there is a change of the system, you simply send the updated data (tasks) to the Davic generic microservices, and they are updated on-the-fly. 

**Davic** framework allows you to transit **computations** in networks. Cool? (Nah... I think it is cool.)
Here key is simple: write scripts in JSON format. So, your scripts are data now. 
I should say this at the beginning: program/script is data, at the first place. 

Here are some example usages of **Davic** (more are coming) 
* [REST-mocker](https://github.com/wfchiang/davic-sample-apps/tree/master/rest-mocker). A mock RESTful service which can be defined "what to mock" **at runtime**. 

Here is a tool set that helps you write scripts in JSON: [Davic-helper](https://github.com/wfchiang/davic-helpers). 