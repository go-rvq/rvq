package getting_started

import (
	. "github.com/go-rvq/rvq/admin/docs/docgo"
)

var OneMinuteQuickStart = Doc(
	Markdown(`
This brief tutorial aims to give you a rapid taste of RVQ's capabilities in the shortest possible time. One standout feature of RVQ is its "presets" module, which swiftly generates [fully functional admin interfaces](/samples/presets-detail-page-cards/customers) like those you see below.

To get started right away:


1. **Install the Command Line Tool**: Run the following command to install the latest version of the RVQ CLI tool:

~~~
$ go install github.com/go-rvq/rvq/admin/docs/cmd/rvq@latest
~~~

2. **Launch RVQ**: Execute the rvq command:

~~~
$ rvq
~~~

You'll be prompted to enter a Go package name. 

And then there are these template options,

- Admin: Content Management System
- Website: Content Management System with Website Examples
- Bare: Simplest Workable Web App

Here we select Admin, The tool will then create an admin app within your current working directory.

3. **Set Up the Database**: Navigate to the newly created package directory and start the database using Docker Compose:

~~~
$ cd <your_package_name>
$ docker-compose up
~~~
This command launches the necessary database services

4. **Run the Admin App**: Open a new terminal window and execute the following commands to load the development environment variables and run the admin application:

~~~
$ source dev_env
$ go run main.go
~~~

With these quick steps, you'll have a fully operational RVQ admin interface up and running, showcasing the remarkable speed and efficiency at which RVQ empowers you to build sophisticated web applications. Explore the interface to witness firsthand the extent of RVQ's power and versatility, all within just one minute!

`),
).Title("1 Minute Quick Start").
	Slug("getting-started/one-minute-quick-start")
