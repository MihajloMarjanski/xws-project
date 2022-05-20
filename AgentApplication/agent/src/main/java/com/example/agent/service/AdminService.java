package com.example.agent.service;

import com.example.agent.model.*;
import com.example.agent.repository.AdminRepository;
import com.example.agent.repository.CompanyOwnerRepository;
import com.example.agent.repository.CompanyRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Service;

import java.util.*;

@Service
public class AdminService {

    @Autowired
    private RoleService roleService;
    @Autowired
    private AdminRepository adminRepository;
    @Autowired
    private CompanyRepository companyRepository;
    @Autowired
    private CompanyOwnerRepository companyOwnerRepository;


    public ResponseEntity<?> approveCompany(Integer id) {
        Optional<Company> company = companyRepository.findById(id);
        if (!company.isPresent())
            return new ResponseEntity<>(HttpStatus.BAD_REQUEST);

        Role role = roleService.findByName("ROLE_COMPANY_OWNER");
        CompanyOwner owner = company.get().getCompanyOwner();
        Set<Role> ownerRoles = owner.getRoles();
        ownerRoles.add(role);
        owner.setRoles(ownerRoles);
        companyOwnerRepository.save(owner);

        company.get().setApproved(true);
        companyRepository.save(company.get());
        return new ResponseEntity<>(HttpStatus.OK);
    }

    public Collection<String> findAllUsernames() {
        return adminRepository.findAllUsernames();
    }

    public boolean isPinOk(String username, Integer pin) {
        Admin user = adminRepository.findByUsername(username);
        if (user == null)
            return false;
        return user.getPin().equals(pin);
    }

    public Admin findByUsername(String username) {
        return adminRepository.findByUsername(username);
    }

    public ResponseEntity<?> getByUsername(String username) {
        if (findByUsername(username) == null)
            return new ResponseEntity<>(HttpStatus.BAD_REQUEST);
        return new ResponseEntity<>(findByUsername(username), HttpStatus.OK);
    }
}
