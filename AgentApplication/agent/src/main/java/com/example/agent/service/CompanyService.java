package com.example.agent.service;

import com.example.agent.model.Company;
import com.example.agent.model.CompanyOwner;
import com.example.agent.model.Role;
import com.example.agent.repository.CompanyOwnerRepository;
import com.example.agent.repository.CompanyRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.dao.DataIntegrityViolationException;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;

import java.sql.SQLIntegrityConstraintViolationException;
import java.util.Collection;
import java.util.List;
import java.util.Optional;

@Service
public class CompanyService {

    @Autowired
    private CompanyOwnerRepository companyOwnerRepository;
    @Autowired
    private CompanyRepository companyRepository;
    @Autowired
    private RoleService roleService;

    public ResponseEntity<?> saveCompanyOwner(CompanyOwner companyOwner) {
        PasswordEncoder passwordEncoder = new BCryptPasswordEncoder();
        try {
            companyOwner.setPassword(passwordEncoder.encode(companyOwner.getPassword()));
            List<Role> roles = roleService.findByName("ROLE_REGISTERED");
            companyOwner.setRole(roles.get(0));
            companyOwnerRepository.save(companyOwner);
            return new ResponseEntity<>(HttpStatus.OK);
        } catch (DataIntegrityViolationException e) {
            e.printStackTrace();
            return new ResponseEntity<>("Already exist user with same username or email", HttpStatus.BAD_REQUEST);
        }
    }

    public ResponseEntity<?> sendCompanyRegistrationRequest(Company company) {
        Optional<CompanyOwner> owner = companyOwnerRepository.findById(company.getCompanyOwner().getId());
        if (!owner.isPresent())
            return new ResponseEntity<>(HttpStatus.BAD_REQUEST);
        else if (companyRepository.findByCompanyOwnerId(owner.get().getId()) != null)
            return new ResponseEntity<>(HttpStatus.BAD_REQUEST);
        companyRepository.save(company);
        return new ResponseEntity<>(HttpStatus.OK);
    }

    public ResponseEntity<?> getOwner(Integer id) {
        Optional<CompanyOwner> owner = companyOwnerRepository.findById(id);
        if(!owner.isPresent())
            return new ResponseEntity<>(HttpStatus.BAD_REQUEST);
        return new ResponseEntity<>(owner.get(), HttpStatus.OK);
    }

    public Collection<String> findAllUsernames() {
        return companyOwnerRepository.findAllUsernames();
    }
}
