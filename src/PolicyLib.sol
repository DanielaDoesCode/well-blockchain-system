// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.8.0 <0.9.0;

import "./Types.sol";

//Static policy evaluation - all logic
library PolicyLib {

    // Evaluate provider's attributes
    function evaluate(
        Types.ProviderAttributes memory caller,
        bytes32 requiredDepartment,
        bytes32 requiredInstitution,
        uint48 evaluatedAt
    ) internal pure returns (bool) {
        if (caller.validUntil < evaluatedAt)       
            return false;
        if (requiredDepartment != bytes32(0) && 
            caller.department  != requiredDepartment) 
            return false;
        if (requiredInstitution != bytes32(0) &&
            caller.institution  != requiredInstitution) 
            return false;
        return true;
    }
}