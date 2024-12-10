# Basic Multiregion Azure Web App
```mermaid
architecture-beta
    group shared(cloud)[Shared Services]
        service kv(database)[Key Vault] in shared
        service acr(database)[Premium ACR] in shared
        service tfstate(disk)[State Storage] in shared

    group primary(cloud)[East US Region]
        service east_vnet(internet)[VNet East] in primary
        service east_app(server)[Web App] in primary
        service east_pe(disk)[Private EP] in primary

    group secondary(cloud)[Central US Region]
        service central_vnet(internet)[VNet Central] in secondary
        service central_app(server)[Web App] in secondary
        service central_pe(disk)[Private EP] in secondary

    group cicd(cloud)[CI/CD]
        service repo(server)[GitHub] in cicd
        service actions(server)[Actions] in cicd

    service tm(internet)[Traffic Manager]

    tm:B -- T:east_app
    tm:B -- T:central_app

    repo:R -- L:actions
    actions:R -- L:acr

    east_vnet:R -- L:east_app
    east_app:B -- T:east_pe

    central_vnet:R -- L:central_app
    central_app:B -- T:central_pe

    kv:T -- B:east_app
    kv:T -- B:central_app

    acr:L -- R:east_app
    acr:R -- L:central_app
```

```mermaid
graph TB
    TM[Azure Traffic Manager] --> EUS[East US Region]
    TM --> CUS[Central US Region]
    
    subgraph "East US"
    EUS --> EUSVNET[VNet East US]
    EUSVNET --> EUSWEBAPP[Web App East US]
    EUSWEBAPP --> EUSPE[Private Endpoint]
    end
    
    subgraph "Central US"
    CUS --> CUSVNET[VNet Central US]
    CUSVNET --> CUSWEBAPP[Web App Central US]
    CUSWEBAPP --> CUSPE[Private Endpoint]
    end
    
    subgraph "Shared Services"
    ACR[Azure Container Registry<br>Basic Tier] --> EUSWEBAPP
    ACR --> CUSWEBAPP
    KV[Key Vault] --> EUSWEBAPP
    KV --> CUSWEBAPP
    end
    
    subgraph "CI/CD"
    GH[GitHub Repository] --> GHA[GitHub Actions]
    GHA --> ACR
    end

    classDef azure fill:#0072C6,stroke:#fff,stroke-width:2px,color:#fff
    classDef network fill:#008272,stroke:#fff,stroke-width:2px,color:#fff
    classDef security fill:#CA2171,stroke:#fff,stroke-width:2px,color:#fff
    
    class TM,EUS,CUS,EUSWEBAPP,CUSWEBAPP,ACR,KV azure
    class EUSVNET,CUSVNET,EUSPE,CUSPE network
    class GH,GHA security
```