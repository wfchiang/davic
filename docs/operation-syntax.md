# Operation Syntax
An operation is represented as an array. 
* The first element is a special mark (a string) **-opt-**. 
* The second element is the operator mark. E.g., **-store-read-** is the mark of the "store-read" operator. 
* The remaining elements are the parameters of the operation. 


## Environment Access Operations

### Store Read
Read value from a store entry. 
#### Operator Mark **-store-read-**
#### Parameters 
1. A single string, the key of the store entry

### Store Write
Write value to a store entry. 
#### Operator Mark **-store-write-**
#### Parameters
1. A single string, the key of the store entry 
2. An expression, the value to put into the store 

### Stack Read 
Read the top value of the stack. 
#### Operator Mark **-stack-read-** 
#### No Parameter


## Function Operations 

### Lambda
A function. 
#### Operator Mark **-lambda-** 
#### Parameters 
1. An expression, the function body 

### Function call 
Call a function with a parameter. 
#### Operator Mark **-fcall-**
#### Parameters 
1. A lambda 
2. An expression, the parameter for calling the lambda 


## Relation Operations 

### Equal 
Check if the left-hand-side valule is equal to the right-hand-side value. 
When doing the comparison, the type of the both side must be equal. 
If not, an error will be raised.  
#### Operator Mark **-rel-eq-**
#### Parameters 
1. Left-hand-side value 
2. Right-hand-side value


## Arithmetic Operations 

### Add
Calculate the summasion of the parameters. 
#### Operator Mark **-arith-add-** 
#### Parameters 
1. Left-hand-side value
2. Right-hand-side value


## Array Operations 


## Object Access Operations 