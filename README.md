# Authentication Service

## Login

- path: `/auth/login`
- method : `POST`
- params:
  - body:

    ```js
    {
        "email":String,
        "password":String
    }
    ```

- returns:

    ```js
    {
        "token": String,
        "success": Number,
        "message": String
    }
    ```

## Signup

- path: `auth/signup`
- method: `POST`
- params:
  - body:

    ```js
    {
        "firstname":String,
        "lastname" : String,
        "password" : String,
        "email" : String,
    }
    ```

- returns:

    ```js
    {
        "success": Number,
        "message": String
    }
    ```
