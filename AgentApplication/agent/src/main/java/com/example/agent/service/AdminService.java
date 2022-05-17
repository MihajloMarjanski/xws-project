package com.example.agent.service;

import com.example.agent.model.Company;
import com.example.agent.model.CompanyOwner;
import com.example.agent.model.Role;
import com.example.agent.repository.AdminRepository;
import com.example.agent.repository.CompanyOwnerRepository;
import com.example.agent.repository.CompanyRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Service;

import java.util.Collection;
import java.util.List;
import java.util.Optional;

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

        List<Role> roles = roleService.findByName("ROLE_COMPANY_OWNER");
        CompanyOwner owner = company.get().getCompanyOwner();
        owner.setRole(roles.get(0));
        companyOwnerRepository.save(owner);

        company.get().setApproved(true);
        companyRepository.save(company.get());
        return new ResponseEntity<>(HttpStatus.OK);
    }

    public Collection<String> findAllUsernames() {
        return adminRepository.findAllUsernames();
    }
}
