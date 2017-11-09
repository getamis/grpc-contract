pragma solidity ^0.4.15;

contract NameService {
    string name;
    event NameSet(string _name);

    function setName(string _name) {
        name = _name;
        NameSet(name);
    }

    function getName() constant returns (string) {
        return name;
    }
}