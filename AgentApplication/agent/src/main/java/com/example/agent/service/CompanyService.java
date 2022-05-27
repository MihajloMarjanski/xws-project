package com.example.agent.service;

import com.example.agent.mapper.CompanyOwnerAdapter;
import com.example.agent.model.*;
import com.example.agent.model.dto.*;
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

import java.util.*;

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

    public ResponseEntity<?> createCompanyOwner(UserDto dto) {
        PasswordEncoder passwordEncoder = new BCryptPasswordEncoder();
        try {
            CompanyOwner companyOwner = new CompanyOwner(dto);
            companyOwner.setSalt(RandomStringInitializer.generateAlphaNumericString(10));
            companyOwner.setPassword(passwordEncoder.encode(companyOwner.getPassword().concat(companyOwner.getSalt())));
            String pin = RandomStringInitializer.generatePin();
            companyOwner.setPin(passwordEncoder.encode(pin.concat(companyOwner.getSalt())));
            companyOwner.setActivated(false);
            companyOwner.setForgotten(0);
            companyOwner.setMissedPasswordCounter(0);
            Role role = roleService.findByName("ROLE_POTENTIAL_OWNER");
            Set<Role> ownerRoles = companyOwner.getRoles();
            ownerRoles.add(role);
            companyOwner.setRoles(ownerRoles);
            companyOwnerRepository.save(companyOwner);
            emailService.sendActivationMailOwnerAsync(findByUsername(companyOwner.getUsername()));
            emailService.sendPin(companyOwner.getEmail(), pin);
            return new ResponseEntity<>(HttpStatus.OK);
        } catch (DataIntegrityViolationException e) {
            e.printStackTrace();
            return new ResponseEntity<>("Already exist user with same username or email", HttpStatus.BAD_REQUEST);
        }
    }

    public ResponseEntity<?> sendCompanyRegistrationRequest(Company company, String ownerUsername) {
        CompanyOwner owner = companyOwnerRepository.findByUsername(ownerUsername);
        if (owner == null)
            return new ResponseEntity<>("Owner with that username does not exist.", HttpStatus.BAD_REQUEST);
        else if (companyRepository.findByCompanyOwnerId(owner.getId()) != null)
            return new ResponseEntity<>("Already have company.", HttpStatus.BAD_REQUEST);
        company.setCompanyOwner(owner);
        companyRepository.save(company);
        for (JobPosition job : company.getPositions()) {
            job.setCompany(findByOwner(owner));
            jobPositionRepository.save(job);
        }
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
        restTemplate.postForObject("http://localhost:8000/jobs/offer", jobOffer, Void.class);
        return new ResponseEntity<>(HttpStatus.OK);
    }

    public CompanyOwner findByOwnerEmail(String email) {
        return companyOwnerRepository.findByEmail(email);
    }

    public void saveOwner(CompanyOwner owner) {
        companyOwnerRepository.save(owner);
    }

    public ResponseEntity<?> sendNewPassword(CompanyOwner owner) {
        PasswordEncoder passwordEncoder = new BCryptPasswordEncoder();
        String newPass = RandomStringInitializer.generateAlphaNumericString(10);
        owner.setPassword(passwordEncoder.encode(newPass.concat(owner.getSalt())));
        owner.setPin(RandomStringInitializer.generatePin());
        owner.setForgotten(1);
        saveOwner(owner);
        emailService.sendNewPassword(owner.getEmail(), newPass);
        emailService.sendPin(owner.getEmail(), owner.getPin());
        return new ResponseEntity<>(HttpStatus.OK);
    }

    public boolean isPinOk(String username, String pin) {
        CompanyOwner user = companyOwnerRepository.findByUsername(username);
        if (user == null)
            return false;
        PasswordEncoder passwordEncoder = new BCryptPasswordEncoder();
        String saltedPin = pin.concat(user.getSalt());
        boolean match = passwordEncoder.matches(saltedPin, user.getPin());
        return match;
    }

    public CompanyOwner findByUsername(String username) {
        return companyOwnerRepository.findByUsername(username);
    }

    public ResponseEntity<?> updateCompanyOwner(OwnerWithCompany companyOwner) {
        CompanyOwner owner = findByUsername(companyOwner.getUsername());
        owner.setEmail(companyOwner.getEmail());
        owner.setFirstName(companyOwner.getFirstName());
        owner.setLastName(companyOwner.getLastName());
        if(!owner.getPassword().equals(companyOwner.getPassword())) {
            PasswordEncoder passwordEncoder = new BCryptPasswordEncoder();
            owner.setPassword(passwordEncoder.encode(companyOwner.getPassword().concat(owner.getSalt())));
            owner.setForgotten(0);
        }
        saveOwner(owner);

        Company company = findByOwner(owner);
        if(company != null) {
            company.setInfo(companyOwner.getCompany().getInfo());
            company.setName(companyOwner.getCompany().getName());
            company.setCity(companyOwner.getCompany().getCity());
            company.setCountry(companyOwner.getCompany().getCountry());

            for (JobPosition job : companyOwner.getCompany().getPositions())
                job.setCompany(company);
            company.setPositions(companyOwner.getCompany().getPositions());
            companyRepository.save(company);
        }
        return new ResponseEntity<>(owner, HttpStatus.OK);
    }

    public ResponseEntity<?> getAll() {
        List<CompanyDto> dtos = new ArrayList<>();
        for(Company company : companyRepository.findAll())
            dtos.add(new CompanyDto(company));
        return new ResponseEntity<>(dtos, HttpStatus.OK);
    }

    public ResponseEntity<?> getOwnerByUsername(String username) {
        CompanyOwner owner = findByUsername(username);
        if (owner == null)
            return new ResponseEntity<>(HttpStatus.BAD_REQUEST);
        return new ResponseEntity<>(CompanyOwnerAdapter.convertOwnerToDto(owner, findByOwner(owner)), HttpStatus.OK);
    }

    private Company findByOwner(CompanyOwner owner) {
        return companyRepository.findByCompanyOwnerId(owner.getId());
    }

    public ResponseEntity<?> getAllApproved() {
        List<Company> dtos = new ArrayList<>();
        for(Company company : companyRepository.findAll())
            if (company.isApproved()) {
                for (Comment comment : company.getComments())
                    comment.setCompany(null);
                for (JobPosition job : company.getPositions()) {
                    for (InterviewInformation info : job.getInterviewInformations())
                        info.setJobPosition(null);
                    job.setCompany(null);
                }
                dtos.add(company);
            }

        return new ResponseEntity<>(dtos, HttpStatus.OK);
    }

    public ResponseEntity<?> getByOwner(String username) {
        Company company = companyRepository.findByCompanyOwnerUsername(username);
        for (Comment comment : company.getComments())
            comment.setCompany(null);
        for (JobPosition job : company.getPositions()) {
            for (InterviewInformation info : job.getInterviewInformations())
                info.setJobPosition(null);
            job.setCompany(null);
        }
        return new ResponseEntity<>(company, HttpStatus.OK);
    }

    public ResponseEntity<?> getApiKey(String username, String password) {
        RestTemplate restTemplate = new RestTemplate();
//        ResponseEntity<ApiKeyDto> keyRes = restTemplate.getForEntity("http://localhost:8000/user/apiKey/" + username + "/" + password, ApiKeyDto.class);
//        ApiKeyDto res = restTemplate.getForObject("http://localhost:8000/user/apiKey/" + username + "/" + password, ApiKeyDto.class);
//        if (res.getApiKey() == "")
//            return new ResponseEntity<>(HttpStatus.BAD_REQUEST);
        return restTemplate.getForEntity("http://localhost:8000/user/apiKey/" + username + "/" + password, ApiKeyDto.class);
    }
}
