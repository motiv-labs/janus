# Quickstart

**Before you start:** Make sure you've [installed Janus](../install/README.md)

In this section, you'll learn how to manage your Janus instance. First we'll
have you start Janus giving in order to give you access to the RESTful Admin
interface, through which you manage your APIs, consumers, and more. Data sent
through the Admin API is stored in Janus's [datastore](../datastore/README.md) (Janus
supports File System and MongoDB).

1. ### Start Janus.

    We recommend using some sort of process control, like [runit](http://smarden.org/runit).
    
    Issue the following command to [start][CLI] Janus:

    ```bash
    $ sv start janus
    ```

    **Note:** The CLI also accepts a configuration (`-c <path_to_config>`)
    option allowing you to point to different configurations.

2. ### Verify that Janus has started successfully

    Once these have finished you should see a message (`Janus started`)
    informing you that Janus is running.

    By default Janus listens on the following ports:

- `:8080` on which Janus listens for incoming HTTP traffic from your
  clients, and forwards it to your upstream services.
- `:8443` on which Janus listens for incoming HTTPS traffic. This port has a
  similar behavior as the `:8080` port, except that it expects HTTPS
  traffic only. This port can be disabled via the configuration file.
- `:8081` on which the Admin API used to configure Janus listens.
- `:8444` on which the Admin API listens for HTTPS traffic.

3. ### Stop Janus.

    As needed you can stop the Janus process by issuing the following
    [command][CLI]:

    ```bash
    $ sv stop janus
    ```

4. ### Reload Janus.

    Issue the following command to [reload][CLI] Janus without downtime:

    ```bash
    $ sv reload janus
    ```
