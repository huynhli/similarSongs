# similarSongs

See it working [here](https://similarsongs.pages.dev)!

Have you ever noticed that after listening to Spotify for a while, you get recommended the same songs every time? Whether it be albums, artists, songs, or playlists, Spotify uses [collaborative filtering](https://www.techaheadcorp.com/blog/spotify-recommendation-system/#:~:text=Instead%20of%20just%20individual%20user,deeper%20understanding%20of%20user%20preferences.) to recommend songs to users. This is summarized as Spotify looking for users who like the same songs as you, 
and recommending songs that *they* like. This *can* be great, but it can also lead to you hearing the exact same song every time you finish an album or playlist. So, that led me to making this.  


A link to our Jira board can be found here: [similarSongs Jira Board](https://csprojectss.atlassian.net/jira/software/projects/SS/boards/4?filter=&groupBy=none)
  
  
To run it locally, clone and run:  
1. the backend and frontend manually. This requires two terminals.  
   - terminal 1:  
       - cd backend  
       - go mod download  
       - go run main.go  
   - terminal 2:  
       - cd frontend  
       - npm i  
       - npm run dev  
3. it using the dockerfiles/docker compose.  
   - dockerfiles would require 1 terminal in each backend/frontend dir running:  
       - docker build -t <YOUR_IMAGE_NAME> .
       - docker run --env-file .env -p 8080:8080 <YOUR_IMAGE_NAME>
   - docker-compose up --build  
       - this requires docker desktop    
5. it with make commands.
   - one terminal with docker desktop
       - make dev  
  
Your backend env file should have the following vars:  
DB_USERNAME=  
DB_PASSWORD=  
DB_HOST=  
DB_PORT=  
DB_NAME=  
JWT_SECRET=  
  
Your frontend env file should have the following vars:  
VITE_API_BASE_URL=  

