## client library

- A class that receives a circuit object and it creates all the metadata and qasm version of the circuit
- If the circuit is not currently recognized (it's not from qiskit, cirq or pennylane), there's another class for custom providers, that requires a bunch of input data from
- backend instances are predefined in the configs (default 3)

## CONSTRAINTS

- must set a config for the max of simultaneous simulators (.toml)
- each simulator node has a queue associated with with the sequence of jobs

## API

- crate job

    - receives: QASM data + metadata
    - via: RPC+protobuf

    - data definition: 
        - metadata: 
        ```js
        {
            n_qubits: Int,
            framework: String,
            submission_date: Datetime,
            depth: Int,
            result_type: String[] (counts, quasi dist, etc.),
            extra: String (JSON stringify),
            target: String (target simulator)
        }
        ```

    - process:
        1. job is registered on history database (ids are generated directly on the sever not in the database)
        2. job is submitted to the rabbitmq queue
        3. rabbitmq assigns the job to the least busy worker
        4. ends
    
    - returns: job uuid

    - after: after a worker gets a job, it must run it in the given backend/provider and at the end save the results in the database. It also needs to update values like finish_time, status, etc.

- get job
    
    - receives: job id
    - via: RPC+protobuf

    - process:
        1. gets the id
        2. retrieves the data from database
        3. pack into protobuf and return the data
        4. ends

    - returns: results data (counts, quasi dist)


- cancel job
    //TODO


## JOBS

- after finishing a job your must delete the qasm file and set the status to `finished` or `failed`
- you can also cancel a job, so after that the worker must delete the qasm file and set the status to `canceled` 
- a same job, can raise different types of results

- statuses: [ pending, running, finished, canceled,  failed]  
- result types: [quasi, counts, expval]