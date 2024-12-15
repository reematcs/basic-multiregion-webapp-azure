# Basic Multiregion Azure Web App

![image info](images/multi-region_azure_web_application.png)

# Multi-Region Azure Health Dashboard

## Overview
The Multi-Region Azure Health Dashboard is a monitoring and failover management application deployed across West US and Central US regions.

## Description
The system uses a three-tier architecture, providing real-time health monitoring and controlled failover operations through Traffic Manager routing, secured by private endpoints and regional VNets. The system centralizes container and secret management while using GitHub Actions for consistent cross-region deployments.

## Architecture Details

### Tiers
1. Frontend/Edge Tier
   - Traffic Manager for global routing
   - Private Endpoints in both regions
   - Health status and failover control interface

2. Application Tier
   - Web Apps in West US (Primary) and Central US (Secondary)
   - Regional VNets (10.1.0.0/16 and 10.2.0.0/16)
   - Health monitoring and failover logic

3. Shared Services Tier
   - Premium Azure Container Registry (ACR)
   - Key Vault for secrets management
   - Terraform State storage
   - GitHub repository and Actions

### Security
- Private endpoint connections
- Isolated virtual networks
- Centralized secrets management
- Secure container registry

### Deployment
- GitHub Actions pipeline
- Consistent cross-region deployment
- Infrastructure as Code via Terraform
- Container-based application deployment

## Application Details

### Core Features
1. Region Information
   - Region name display
   - Instance hostname
   - Request timestamp
   - Container version tracking

2. Health Monitoring
   - Traffic Manager endpoint status
   - Key Vault connectivity
   - Container Registry access
   - Regional health status

3. Failover Management
   - Primary/Secondary status indicators
   - Manual failover controls
   - Failover history tracking

### API Endpoints
1. System Information
   - `/api/system` - Basic region and instance information
   - `/api/health/status` - Comprehensive health status
   - `/api/health/live` - Traffic Manager health probe

2. Failover Control
   - `/api/failover/trigger` - Manual failover initiation
   - `/api/failover/history` - Failover history and status

### Monitoring Capabilities
- Real-time health status
- Service connectivity monitoring
- Container deployment tracking
- Cross-region status comparison