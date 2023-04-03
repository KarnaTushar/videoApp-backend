# videoApp-backend

follwing are the steps for installation
1. clone this repo in dev environment
2. fill config.json values(you can refer the existing values or this -> https://appbuilder-docs.agora.io/turn-key/guides/Backend/Credentials)
3. Make sure you have mysql and go installed on your system
4. Once filling the entire details you can run this code either through docker(explained in step 5) or by serving the build file
5. through docker: fill docker-yaml file or create your own and then compose it
6. through build file: create a build file using command 'go build -o ./server ./cmd/video_conferencing'
7. Since migration flag is set to false you have to create database and tables within it
8. refer commands from ./migrations/migrations folder to create or drop sql tables using .up.sql or .down.sql file respectively
9. this server will run on port 8080
