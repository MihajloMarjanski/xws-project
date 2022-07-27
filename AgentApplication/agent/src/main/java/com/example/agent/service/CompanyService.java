package com.example.agent.service;

import com.example.agent.mapper.CompanyOwnerAdapter;
import com.example.agent.model.*;
import com.example.agent.model.dto.*;
import com.example.agent.repository.CompanyOwnerRepository;
import com.example.agent.repository.CompanyRepository;
import com.example.agent.repository.JobPositionRepository;
import com.example.agent.security.tokenUtils.JwtTokenUtils;
import lombok.extern.slf4j.Slf4j;
import org.apache.http.ssl.SSLContextBuilder;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.autoconfigure.neo4j.Neo4jProperties;
import org.springframework.dao.DataIntegrityViolationException;
import org.springframework.http.*;
import org.springframework.http.client.HttpComponentsClientHttpRequestFactory;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestTemplate;

import javax.net.ssl.SSLContext;
import javax.servlet.http.HttpServletRequest;
import java.util.*;

@Service
@Slf4j
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
    @Autowired
    JwtTokenUtils tokenUtils;

    public ResponseEntity<?> createCompanyOwner(UserDto dto, HttpServletRequest request) {
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
            log.info("Ip: {}, username: {}, Client successfully created!", request.getRemoteAddr(), companyOwner.getUsername());
            return new ResponseEntity<>(HttpStatus.OK);
        } catch (DataIntegrityViolationException e) {
            log.error("Ip: {}, Client not created! Already exist user with same username or email.", request.getRemoteAddr(), e);
            return new ResponseEntity<>("Already exist user with same username or email", HttpStatus.BAD_REQUEST);
        }
    }

    public ResponseEntity<?> sendCompanyRegistrationRequest(Company company, String ownerUsername) {
        CompanyOwner owner = companyOwnerRepository.findByUsername(ownerUsername);
        if (owner == null){
            log.warn("Username: {}, Company owner with that username doesn't exist", ownerUsername);
            return new ResponseEntity<>("Owner with that username does not exist.", HttpStatus.BAD_REQUEST);
        }
        else if (companyRepository.findByCompanyOwnerId(owner.getId()) != null){
            log.warn("Username: {}, Company owner already has company!", ownerUsername);
            return new ResponseEntity<>("Already have company.", HttpStatus.BAD_REQUEST);
        }
        company.setCompanyOwner(owner);
        companyRepository.save(company);
        for (JobPosition job : company.getPositions()) {
            job.setCompany(findByOwner(owner));
            jobPositionRepository.save(job);
        }
        log.info("Username: {} company: {}, Company added to owner successfully!", ownerUsername, company.getName());
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
        restTemplate.postForObject("https://localhost:8000/jobs/offer", jobOffer, Void.class);
        log.info("Job offer sent to dislink");
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
        Calendar c = Calendar.getInstance();
        c.setTime(user.getPinCreatedDate());
        c.add(Calendar.MINUTE, 1);

        if (user.getPin().equals("") || c.getTime().before(new Date())) {
            return false;
        }
        PasswordEncoder passwordEncoder = new BCryptPasswordEncoder();
        String saltedPin = pin.concat(user.getSalt());
        boolean match = passwordEncoder.matches(saltedPin, user.getPin());
        return match;
    }

    public CompanyOwner findByUsername(String username) {
        return companyOwnerRepository.findByUsername(username);
    }

    public ResponseEntity<?> updateCompanyOwner(OwnerWithCompany companyOwner, HttpServletRequest request) {
        CompanyOwner owner = findByUsername(companyOwner.getUsername());
        owner.setEmail(companyOwner.getEmail());
        owner.setFirstName(companyOwner.getFirstName());
        owner.setLastName(companyOwner.getLastName());
        if(!owner.getPassword().equals(companyOwner.getPassword())) {
            PasswordEncoder passwordEncoder = new BCryptPasswordEncoder();
            owner.setPassword(passwordEncoder.encode(companyOwner.getPassword().concat(owner.getSalt())));
            owner.setForgotten(0);
        }
        log.info("Ip: {}, username: {}, Company owner updated successfully!", request.getRemoteAddr(), companyOwner.getUsername());
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
            log.info("Ip: {}, username: {}, company: {}, Company updated by user successfully!", request.getRemoteAddr(), companyOwner.getUsername(), company.getName());
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
        log.info("Username: {}, Api key sent to dislink!", username);
        return restTemplate.getForEntity("https://localhost:8000/user/apiKey/" + username + "/" + password, ApiKeyDto.class);
    }

    public ResponseEntity<?> searchOffers(String text) {
        if (text == null)
            text = "";
        RestTemplate restTemplate = new RestTemplate();

        String jwt = tokenUtils.generateToken("AgentApplication", null);
        HttpHeaders headers = new HttpHeaders();
        headers.setContentType(MediaType.APPLICATION_JSON);
        headers.setBearerAuth(jwt);

        HttpEntity request = new HttpEntity(headers);
        log.info("Username: {},Search offer sent to dislink!");
        return restTemplate.exchange("https://localhost:8000/jobs/search/" + text, HttpMethod.GET, request, String.class);
    }

    public void send2factorAuthPin(CompanyOwner owner) {
        PasswordEncoder passwordEncoder = new BCryptPasswordEncoder();
        String pin = RandomStringInitializer.generatePin();
        owner.setPin(passwordEncoder.encode(pin.concat(owner.getSalt())));
        owner.setPinCreatedDate(new Date());
        companyOwnerRepository.save(owner);
        emailService.send2factorAuthPin(owner.getEmail(), pin);
    }
}
