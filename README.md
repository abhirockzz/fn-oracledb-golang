# Example for Go functions in Fn talking to Oracle DB

## step 0

- `git clone https://github.com/abhirockzz/fn-oracledb-golang`
- run `seed-db.sql` to create the table in your DB
- replace `config` section `func.yaml` with DB connectivity details for your specific environment

## step 1

- `fn start`

From a different terminal

if you want to test locally ignore the next two steps

- `docker login` (use your docker registry credentials)
- `export FN_REGISTRY=<name of your docker repository>`

> your function docker image name will end up being - `<docker repo name>/<function name in func.yaml>:<version in func.yaml>`

Moving on....

- `cd fn-oracledb-golang`
- `fn -v deploy --all --local` (`-v` will activate verbose mode)

> `--local` will build & push docker images locally (and run it from there). Remove it if you want use a dedicated/external Docker registry

All your functions (create, read, update, delete) should now be deployed. Check it using `fn inspect app fn-oradb-go-app` and `fn list routes --app fn-oradb-go-app`

## Behind the scenes

### Build and runtime Docker images

Custom Docker images have been used for `build_image` and `run_image` (see `func.yaml`) - this Docker image pre-packages a [Go Oracle driver](https://github.com/mattn/go-oci8). You can also build you own image

### Hot functions and DB connection

DB connection is not closed at end of every function invocation - this is to leverage hot functions. If the function is called before the `idle_timeout` (default is 30 seconds and can be configured via `func.yaml`), the call will be routed to a hot function (already used) DB connection will not be repeated again 

### app.yaml

Bare minimum config file which defines the name of the app i.e. `fn-oradb-go-app`

## Test

you can test in two ways

- Fn CLI using `fn call`
- `curl` the endpoint of your functions as per `fn list routes fn-oradb-go-app`


### Create

`echo '{"emp_email": "a@b.com","emp_name": "abhishek","emp_dept": "Product Divison"}' | fn call fn-oradb-go-app /create`

Create as many as you want

### Read

- `fn call fn-oradb-go-app /read` (to fetch all employees)
- `echo a@b.com | fn call fn-oradb-go-app /read` (to fetch employee with email `a@b.com`)

### Update

It is possible to update the department of an employee

`echo -d '{"emp_email": "a@b.com", "emp_dept": "Support Operations"}' | fn call fn-oradb-go-app /update`

> check to make sure - `echo a@b.com | fn call fn-oradb-go-app /read`

### Delete

Use employee email to specify which employee record you want to delete

`echo a@b.com | fn call fn-oradb-go-app /delete`

> check to make sure - `echo a@b.com | fn call fn-oradb-go-app /read`