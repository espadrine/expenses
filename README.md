# Expense Management

*Perform statistics and predict your needs.*

## CLI Interface

### Initialization

When running the program for the first, time, it will initialize the database
in ~/.config/expense/db.sqlite and initialize the default user (using `$USER`).

You may explicitly create and manage users with the following commands:

    $ expense user create my_user_name
    6gqsxomsj3hir4msckf4wxtmsuoqa3w5
    $ expense user name 6gqsxomsj3hir4msckf4wxtmsuoqa3w5
    my_user_name
    $ expense user id my_user_name
    6gqsxomsj3hir4msckf4wxtmsuoqa3w5
    $ expense user list
    6gqsxomsj3hir4msckf4wxtmsuoqa3w5	my_user_name

### Entering expenses

You may add entries one by one, receiving their generated ID,
or import them bulk:

    $ expense entry add date=2023-09-06 amount=6.03 currency=EUR tag=food tag=cheese label='Parmigiano AOP 1Kg'  # Not yet available
    1
    $ expense entry show 1  # Not yet available
    1	2023-09-06	6.03	EUR	food, cheese	Parmigiano AOP 1Kg
    $ expense entry export  # Not yet available
    id	date	amount	currency	tags	label
    1	2023-09-06	6.03	EUR	food, cheese	Parmigiano AOP 1Kg
    $ expense entry import filetype=csv <tx.csv  # Not yet available

You can then browse statistics on them, or add them from the browser:

    $ expense ui  # Not yet available
    https://localhost:8576/
