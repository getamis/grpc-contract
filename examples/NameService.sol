pragma solidity ^0.4.15;

contract NameService {
    string name;
    event NameSet(string _name);

    function setName(string _name) public {
        name = _name;
        NameSet(name);
    }

    function getName() public constant returns (string) {
        return name;
    }
}