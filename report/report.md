## ${resource_type}@${api_version} - ${error_code_in_title}

### Description

I found differences between PUT request body and GET response:

${diff_description}

```json
${diff_json}
```

### Details

1. ARM Fully-Qualified Resource Type
```
${resource_type}
```

2. API Version
```
${api_version}
```

3. Swagger issue type
```
Swagger Correctness
```

4. OperationId
```
${operation_id}
```

5. Swagger GitHub permalink
```
TODO, 
e.g., https://github.com/Azure/azure-rest-api-specs/blob/60723d13309c8f8060d020a7f3dd9d6e380f0bbd
/specification/compute/resource-manager/Microsoft.Compute/stable/2020-06-01/compute.json#L9065-L9101
```

6. Error code
```
${error_code_in_block}
```

7. Request traces
```
${request_traces}
```

### Links
1. [Semantic and Model Violations Reference](https://github.com/Azure/azure-rest-api-specs/blob/main/documentation/Semantic-and-Model-Violations-Reference.md)
2. [S360 action item generator for Swagger issues](https://aka.ms/swaggers360)