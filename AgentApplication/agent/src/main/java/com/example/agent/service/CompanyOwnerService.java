package com.example.agent.service;

import com.example.agent.model.CompanyOwner;
import com.example.agent.model.Role;
import com.example.agent.repository.CompanyOwnerRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.List;

@Service
public class CompanyOwnerService {

    @Autowired
    private RoleService roleService;
    @Autowired
    private CompanyOwnerRepository companyOwnerRepository;

    public void save(CompanyOwner companyOwner) {
        List<Role> roles = roleService.findByName("ROLE_COMPANY_OWNER");
        companyOwner.setRole(roles.get(0));
        companyOwnerRepository.save(companyOwner);
    }
}
