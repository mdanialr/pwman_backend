# Password Manager API
Simple password manager service

## How to Use
1. Download the binary file from [GitHub Releases](https://github.com/mdanialr/pwman_backend/releases).
2. Extract the downloaded binary file then make sure it's executable.
    ```bash
    tar -xzf pwman_backend....tar.gz
    chmod u+x pwman_backend
    ```
3. Create configuration file from the template.
    ```bash
    cp app.yml.example app.yml
    ```
4. Edit the app config file as needed. You can check the template for explanation of each field.
5. Generate new secret string for TOTP.
    ```bash
    ./pwman_backend -gen
    # will output something like 'Your secret: KBXXXX............'
    ```
6. Put that secret string to `app.yml` in section `cred.secret`.
7. Use 2FA apps such as __Microsoft Authenticator__ or __Authy__ or similar, add new account and manually copy secret string
   above, or you can generate the QR code file using `./sns_backend -qr "/my/path/"` and finally just scan that QR image.
8. Verify that the TOTP is valid
    ```bash
    ./pwman_backend -verify 123456
    # will output 'ERR: INVALID' if the given otp is invalid otherwise should output 'VERIFIED'
    ```
   If already output __VERIFIED__ then you may continue for next step.
9. Run migration and seeder (optional). _migration alone should be sufficient since that will create the tables_.
    ```bash
    ./pwman_backend -migrate -seed
    # if only need migration then just use `-migrate`
    ```
10. Change debug in `app.yml` to `false`, then run the app.
    ```bash
    ./pwman_backend
    ```
11. Check the log file that should be resided in directory that you put in config. There are should be 3 logs file:
  - `app-log` is for fiber log access log, contain all endpoints that has been hit by client.
  - `log` is for internal log, for example if failed to query to db, this app's host and port, etc.
  - `gorm-log` just as the name suggest, GORM-related log file.

### Optional (_Integrate with systemd_)
  ```bash
  [Unit]
  Description=instance to serve pwman api service
  After=network.target

  [Service]
  User=root
  Group=your-username
  ExecStart=/bin/sh -c "cd /path/to/binary/file && ./pwman_backend"

  [Install]
  WantedBy=multi-user.target
  ```
1. Save above _systemd script_ in `/etc/systemd/system/` with a filename maybe something like `pwman_backend.service`.
2. Run and enable systemd, so it will run even after reboot.
  ```bash
  sudo systemctl enable pwman_backend.service --now
  ```

# License
This project is licensed under the **MIT License** - see the [LICENSE](LICENSE "LICENSE") file for details.
