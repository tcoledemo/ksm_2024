name: Build and Sign Go Application

on:
  workflow_dispatch:

jobs:
  build-and-sign:
    environment: prod
    runs-on: ubuntu-latest

    steps:
      - name: Retrieve secrets from Keeper
        id: ksecrets
        uses: Keeper-Security/ksm-action@master
        with:
          keeper-secret-config: ${{ secrets.KSM_CONFIG }}
          secrets: |-
            o5IuHXuXUZyNqxLKB2o1ew/field/password > PASSWORD
            o5IuHXuXUZyNqxLKB2o1ew/file/private.key > file:/tmp/private.key
            o5IuHXuXUZyNqxLKB2o1ew/file/public.key > file:/tmp/public.key
            
      - name: Checkout code
        uses: actions/checkout@v2

      - name: Setup Go
        uses: actions/setup-go@v2
        with:
          go-version: '1.18'

      - name: Import GPG Key
        run: |
          echo "${{ steps.ksecrets.outputs.PASSWORD }}" | gpg --batch --yes --passphrase-fd 0 --import /tmp/private.key
        env:
          GPG_PRIVATE_KEY: ${{ steps.ksecrets.outputs.PASSWORD }}

      - name: Build Go application
        run: go build -o hello-world main.go

      - name: Sign the binary
        run: |
          echo "${{ steps.ksecrets.outputs.PASSWORD }}" | gpg --batch --yes --passphrase-fd 0 --pinentry-mode loopback --armor --output hello-world.asc --detach-sign hello-world
        env:
          GPG_PASSPHRASE: ${{ steps.ksecrets.outputs.PASSWORD }}

      - name: Verify the signature
        run: |
          gpg --verify hello-world.asc hello-world

      - name: Run the binary
        run: ./hello-world
        
      - name: Upload signed binary
        uses: actions/upload-artifact@v2
        with:
          name: signed-binary
          path: hello-world.asc
