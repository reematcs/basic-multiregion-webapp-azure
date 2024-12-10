# Basic Multiregion Azure Web App
```
graph TB
    TM[Azure Traffic Manager] --> EUS[East US Region]
    TM --> WUS[West US Region]
    
    subgraph "East US"
    EUS --> EUSVNET[VNet East US]
    EUSVNET --> EUSWEBAPP[Web App East US]
    EUSWEBAPP --> EUSPE[Private Endpoint]
    end
    
    subgraph "West US"
    WUS --> WUSVNET[VNet West US]
    WUSVNET --> WUSWEBAPP[Web App West US]
    WUSWEBAPP --> WUSPE[Private Endpoint]
    end
    
    subgraph "Shared Services"
    ACR[Azure Container Registry<br>Basic Tier] --> EUSWEBAPP
    ACR --> WUSWEBAPP
    KV[Key Vault] --> EUSWEBAPP
    KV --> WUSWEBAPP
    end
    
    subgraph "CI/CD"
    GH[GitHub Repository] --> GHA[GitHub Actions]
    GHA --> ACR
    end

    classDef azure fill:#0072C6,stroke:#fff,stroke-width:2px,color:#fff
    classDef network fill:#008272,stroke:#fff,stroke-width:2px,color:#fff
    classDef security fill:#CA2171,stroke:#fff,stroke-width:2px,color:#fff
    
    class TM,EUS,WUS,EUSWEBAPP,WUSWEBAPP,ACR,KV azure
    class EUSVNET,WUSVNET,EUSPE,WUSPE network
    class GH,GHA security
```