## client library

- A class that recieves a circuit object and it creates all the metadata and qasm version of the circuit
- If the circuit is not currently recognized (it's not from qiskit, cirq or pennylane), there's another class for custom providers, that requires a bunch of input data from
- backend instances are predefined in the configs (default 3)

## API

- crate job

    - recieves: QASM data + metadata
    - via: RPC+protobuf

    - data definition: 
        - metadata: {
            n_qubits: int,
            framework: String,
            submission_date: datetime,
            depth: int,
            result_type: String (counts, quasi dist, etc.)
            extra: String (JSON stringify)
        }

    - process:
        1. job is registered on history database (ids are generated directly on the sever not in the database)
        2. job is submitted to the rabbitmq queue
        3. rabbitmq assigns the job to the least busy worker
        4. ends
    
    - returns: job uuid

    - after: after a worker gets a job, it must run it in the given backend/provider and at the end save the results in the database. It also needs to update values like finsih_time, status, etc.

- get job
    
    - recieves: job id
    - via: RPC+protobuf

    - process:
        1. gets the id
        2. retrieves the data from database
        3. pack into protobuf and return the data
        4. ends

    - returns: results data (counts, quasi dist)
