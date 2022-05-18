package com.example.agent.service;

import com.example.agent.model.*;
import com.example.agent.model.dto.JobOffer;
import com.example.agent.repository.CompanyOwnerRepository;
import com.example.agent.repository.CompanyRepository;
import com.example.agent.repository.JobPositionRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.dao.DataIntegrityViolationException;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestTemplate;

import java.sql.SQLIntegrityConstraintViolationException;
import java.util.Collection;
import java.util.List;
import java.util.Optional;
import java.util.Set;

@Service
public class CompanyService {

    @Autowired
    private CompanyOwnerRepository companyOwnerRepository;
    @Autowired
    private CompanyRepository companyRepository;
    @Autowired
    private RoleService roleService;
    @Autowired
    private JobPositionRepository jobPositionRepository;
    @Autowired
    EmailService emailService;

    public ResponseEntity<?> saveCompanyOwner(CompanyOwner companyOwner) {
        PasswordEncoder passwordEncoder = new BCryptPasswordEncoder();
        try {
            companyOwner.setSalt(RandomStringInitializer.generateAlphaNumericString(10));
            companyOwner.setPassword(passwordEncoder.encode(companyOwner.getPassword().concat(companyOwner.getSalt())));
            companyOwner.setPin(RandomStringInitializer.generatePin());
            companyOwner.setActivated(false);
            companyOwner.setForgotten(0);
            Role role = roleService.findByName("ROLE_POTENTIAL_OWNER");
            Set<Role> ownerRoles = companyOwner.getRoles();
            ownerRoles.add(role);
            companyOwner.setRoles(ownerRoles);
            companyOwnerRepository.save(companyOwner);
            emailService.sendActivationMailOwnerAsync(findByUsername(companyOwner.getUsername()));
            emailService.sendPin(companyOwner.getEmail(), companyOwner.getPin());
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

    public ResponseEntity<?> getAllJobs(Integer companyId) {
        List<JobPosition> positions = jobPositionRepository.findAllByCompanyId(companyId);
        return new ResponseEntity<>(positions, HttpStatus.OK);
    }

    public ResponseEntity<?> createJobOffer(JobOffer jobOffer) {
        RestTemplate restTemplate = new RestTemplate();
        return restTemplate.postForObject("http://localhost:8000/jobs/offer", jobOffer, ResponseEntity.class);
    }

    public CompanyOwner findByOwnerEmail(String email) {
        return companyOwnerRepository.findByEmail(email);
    }

    public void save(CompanyOwner owner) {
        companyOwnerRepository.save(owner);
    }

    public ResponseEntity<?> sendNewPassword(CompanyOwner owner) {
        owner.setPassword(RandomStringInitializer.generateAlphaNumericString(10));
        owner.setPin(RandomStringInitializer.generatePin());
        owner.setForgotten(1);
        save(owner);
        emailService.sendNewPassword(owner.getEmail(), owner.getPassword());
        emailService.sendPin(owner.getEmail(), owner.getPin());
        return new ResponseEntity<>(HttpStatus.OK);
    }

    public boolean isPinOk(String username, Integer pin) {
        CompanyOwner user = companyOwnerRepository.findByUsername(username);
        if (user == null)
            return false;
        return user.getPin().equals(pin);
    }

    public CompanyOwner findByUsername(String username) {
        return companyOwnerRepository.findByUsername(username);
    }
}
