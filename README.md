# Digivance MVC Application Framework
Welcome to the Digivance MVC Application Framework is an open source Go package intended to provide a lightweight Model View Controller (MVC) application coding library with patterns and conventions similar to more widely used frameworks in other languages (Such as Asp.Net or Rails).

## Version 0.3.0
This version enhances the internal error handling, logging and configuration methods of the package. We have introduced a new ConfigurationManager object that is used in Lieu 
of previous versions Application configuration members. This ConfigurationManager object also allows you to save and load configurations to make is easier to deploy your
application between various environments.

Additionally, we have made some updates to the logging functions adding new formatted wrappers that allow the caller to more easily call methods such as:

```
mvcapp.LogErrorf("Failed to perform task: %v", err)
```

Finally, some updates have been applied that now allow for the caller to use ~/ and ./ to signify the application path when querying the file system. For example,
callers can now set the log file name, relative to the application path as such:

```
mvcapp.SetLogFilename("./my_application.log")
```

## Version 0.2.0
This version enhances the MVC Application Framework with email and content bundling functionality.

## Version 0.1.0
The current master branch is version 0.1.0 of the Digivance MvcApp Framework. This is the initial stable alpha release of the code base. This version contains the basic framework required to design an interactive MVC website / web application.

## References:

> #### Website:
> You can read more about this project and it's maintainers on the MvcApp page at Digivance technologies: [MVC Application Framework](https://digivance.com/services/mvcapp).

> #### Roadmap
> The current road map for development of this package can be found under the [GitHub Issues Board](https://github.com/Digivance/mvcapp/milestones). Note that we use milestones to define major version releases and we use the issues in these milestones to lay out the functionality. 

> #### License
> This framework is developed and released under the Lesser GNU Public Library license. More information can be found in the project [License file](https://github.com/digivance/mvcapp/blob/master/LICENSE).