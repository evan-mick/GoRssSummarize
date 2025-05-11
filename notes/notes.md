

az acr show --name headlinercontcmd --resource-group headliner-cmd-rs -query loginServer --output tsv

docker build -t headlinercontcmd.azurecr.io/headliner:latest
docker push headlinercontcmd.azurecr.io/headliner:latest

the first part is the headliner login server

headlinerconcmd -> headliner container
headliner-cmd-rs -> resource group

headlinercontcmd.azurer.io/headliner:latest

environment: headlinerenv

headlinercontcmd.proudhill-dd259f46.eastus.azurecontainerapps.io

 Of note, typo here with environment being "headlinerenv" without headliner's l.
az containerapp create --name headlinercontcmd --resource-group headliner-cmd-rs --environment headlinerenv --image "headlinercontcmd.azurecr.io/headliner:latest/headliner:latest"  --registry-server "headlinercontcmd.azurecr.io/headliner:latest" --registry-identity system --target-port 8080 --ingress external  
