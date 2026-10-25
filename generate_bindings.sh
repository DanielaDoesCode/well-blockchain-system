### Script to export bindings and create go files to use the bindings

cat out/UserRegistry.sol/UserRegistry.json | jq '.abi' > well-backend/contracts/UserRegistry.abi
cat out/HistoryFactory.sol/HistoryFactory.json | jq '.abi' > well-backend/contracts/HistoryFactory.abi
cat out/Registry.sol/WellRegistry.json | jq '.abi' > well-backend/contracts/Registry.abi
cat out/History.sol/History.json | jq '.abi' > well-backend/contracts/History.abi
cat out/Record.sol/Record.json | jq '.abi' > well-backend/contracts/Record.abi
cat out/EmergencySummary.sol/EmergencySummary.json | jq '.abi' > well-backend/contracts/EmergencySummary.abi
cat out/EmergencySummaryFactory.sol/EmergencySummaryFactory.json | jq '.abi' > well-backend/contracts/EmergencySummaryFactory.abi
cat out/IncidentManager.sol/IncidentManager.json | jq '.abi' > well-backend/contracts/IncidentManager.abi
cat out/RiskAssessment.sol/RiskAssessment.json | jq '.abi' > well-backend/contracts/RiskAssessment.abi
cat out/ZKAuthVerifier.sol/ZKAuthVerifier.json | jq '.abi' > well-backend/contracts/ZKAuthVerifier.abi
cat out/ResearchContributer.sol/ResearchContributer.json | jq '.abi' > well-backend/contracts/ResearchContributer.abi
cat out/ResearchVerifier.sol/ResearchVerifier.json | jq '.abi' > well-backend/contracts/ResearchVerifier.abi
cat out/LogAnchor.sol/LogAnchor.json | jq '.abi' > well-backend/contracts/LogAnchor.abi

abigen --abi well-backend/contracts/UserRegistry.abi --pkg contracts --type UserRegistry --out well-backend/contracts/UserRegistry.go
abigen --abi well-backend/contracts/HistoryFactory.abi --pkg contracts --type HistoryFactory --out well-backend/contracts/HistoryFactory.go
abigen --abi well-backend/contracts/Registry.abi --pkg contracts --type Registry --out well-backend/contracts/Registry.go
abigen --abi well-backend/contracts/History.abi --pkg contracts --type History --out well-backend/contracts/History.go
abigen --abi well-backend/contracts/Record.abi --pkg contracts --type Record --out well-backend/contracts/Record.go
abigen --abi well-backend/contracts/EmergencySummary.abi --pkg contracts --type EmergencySummary --out well-backend/contracts/EmergencySummary.go
abigen --abi well-backend/contracts/EmergencySummaryFactory.abi --pkg contracts --type EmergencySummaryFactory --out well-backend/contracts/EmergencySummaryFactory.go
abigen --abi well-backend/contracts/IncidentManager.abi --pkg contracts --type IncidentManager --out well-backend/contracts/IncidentManager.go
abigen --abi well-backend/contracts/RiskAssessment.abi --pkg contracts --type RiskAssessment --out well-backend/contracts/RiskAssessment.go
abigen --abi well-backend/contracts/ZKAuthVerifier.abi --pkg contracts --type ZKAuthVerifier --out well-backend/contracts/ZKAuthVerifier.go
abigen --abi well-backend/contracts/ResearchContributer.abi --pkg contracts --type ResearchContributer --out well-backend/contracts/ResearchContributer.go
abigen --abi well-backend/contracts/ResearchVerifier.abi --pkg contracts --type ResearchVerifier --out well-backend/contracts/ResearchVerifier.go
abigen --abi well-backend/contracts/LogAnchor.abi --pkg contracts --type LogAnchor --out well-backend/contracts/LogAnchor.go
