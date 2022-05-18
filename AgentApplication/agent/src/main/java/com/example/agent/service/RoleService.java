package com.example.agent.service;

import com.example.agent.model.Role;
import com.example.agent.repository.RoleRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.List;
import java.util.Set;


@Service
public class RoleService {

    @Autowired
    private RoleRepository roleRepository;
    

    public Role findById(Long id) {
        Role auth = this.roleRepository.getOne(id);
        return auth;
    }

    public Role findByName(String name) {
        return this.roleRepository.findByName(name);
    }


}
