FROM launcher.gcr.io/google/aspnetcore:2.1
ADD ./ /app
ENV ASPNETCORE_URLS=http://*:${PORT}
WORKDIR /app
ENTRYPOINT [ "dotnet", "TestApp-2.1.dll" ]
