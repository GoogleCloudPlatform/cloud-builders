FROM gcr.io/google-appengine/aspnetcore:1.1
ADD ./ /app
ENV ASPNETCORE_URLS=http://*:${PORT}
WORKDIR /app
ENTRYPOINT [ "dotnet", "TestApp.dll" ]
