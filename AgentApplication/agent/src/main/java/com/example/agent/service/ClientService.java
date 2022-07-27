package com.example.agent.service;

import com.example.agent.model.*;
import com.example.agent.model.dto.UserDto;
import com.example.agent.repository.*;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.dao.DataIntegrityViolationException;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;

import javax.servlet.http.HttpServletRequest;
import java.io.IOException;
import java.net.URISyntaxException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.*;
import java.util.stream.Collectors;
import java.util.stream.Stream;

@Service
@Slf4j
public class ClientService {

    @Autowired
    ClientRepository clientRepository;
    @Autowired
    RoleService roleService;
    @Autowired
    CommentRepository commentRepository;
    @Autowired
    JobPositionRepository jobPositionRepository;
    @Autowired
    InterviewInformationRepository interviewInformationRepository;
    @Autowired
    EmailService emailService;
    @Autowired
    CompanyRepository companyRepository;

    public ResponseEntity<?> create(UserDto dto, HttpServletRequest request) {
        PasswordEncoder passwordEncoder = new BCryptPasswordEncoder();
        try {
            Client client = new Client(dto);
            client.setSalt(RandomStringInitializer.generateAlphaNumericString(10));
            client.setPassword(passwordEncoder.encode(client.getPassword().concat(client.getSalt())));
            String pin = RandomStringInitializer.generatePin();
            client.setPin(passwordEncoder.encode(pin.concat(client.getSalt())));
            client.setActivated(false);
            client.setForgotten(0);
            client.setMissedPasswordCounter(0);
            Role role = roleService.findByName("ROLE_CLIENT");
            Set<Role> clientRoles = client.getRoles();
            clientRoles.add(role);
            client.setRoles(clientRoles);
            clientRepository.save(client);
            emailService.sendActivationMailClientAsync(findByUsername(client.getUsername()));
            emailService.sendPin(client.getEmail(), pin);
            log.info("Ip: {}, username: {}, Client successfully created!", request.getRemoteAddr(), client.getUsername());
            return new ResponseEntity<>(HttpStatus.OK);
        } catch (DataIntegrityViolationException e) {
            log.error("Ip: {}, Client not created! Already exist user with same username or email.", request.getRemoteAddr(), e);
            return new ResponseEntity<>("Already exist user with same username or email", HttpStatus.BAD_REQUEST);
        }
    }

    public ResponseEntity<?> createComment(Comment comment, Integer companyId) {
        Authentication authentication = SecurityContextHolder.getContext().getAuthentication();
        Company company = companyRepository.getById(companyId);
        comment.setCompany(company);
        commentRepository.save(comment);
        log.info("Username: {}, company: {}, User created comment on company successfully!", authentication.getPrincipal().toString(), company.getName());
        return new ResponseEntity<>(HttpStatus.OK);
    }

    public ResponseEntity<?> updateSalary(JobPosition jobSalary) {
        Optional<JobPosition> job = jobPositionRepository.findById(jobSalary.getId());
        Authentication authentication = SecurityContextHolder.getContext().getAuthentication();
        if (!job.isPresent())
            return new ResponseEntity<>(HttpStatus.BAD_REQUEST);
        Random rand = new Random();
        job.get().setAvgSalary(rand.nextInt(2000 - 1000) + 1000);
        jobPositionRepository.save(job.get());
        log.info("Username: {}, job: {}, Job salary updated successfully!", authentication.getPrincipal().toString(), job.get().getName());
        return new ResponseEntity<>(HttpStatus.OK);
    }

    public ResponseEntity<?> addInterviewInformation(InterviewInformation interviewInformation, Integer jobId) {
        interviewInformation.setJobPosition(jobPositionRepository.getById(jobId));
        Authentication authentication = SecurityContextHolder.getContext().getAuthentication();
        log.info("Username: {}, Interview info updated successfully!", authentication.getPrincipal().toString());
        interviewInformationRepository.save(interviewInformation);
        return new ResponseEntity<>(HttpStatus.OK);
    }

    public void save(Client client) {
        clientRepository.save(client);
    }

    public Client findByEmail(String email) {
        return clientRepository.findByEmail(email);
    }

    public ResponseEntity<?> sendNewPassword(Client client) {
        PasswordEncoder passwordEncoder = new BCryptPasswordEncoder();
        String password = RandomStringInitializer.generateAlphaNumericString(10);
        client.setPassword(passwordEncoder.encode(password.concat(client.getSalt())));
        client.setForgotten(1);
        client.setPin(RandomStringInitializer.generatePin());
        save(client);
        emailService.sendNewPassword(client.getEmail(), password);
        emailService.sendPin(client.getEmail(), client.getPin());
        log.info("Username: {}, New temp password sent to user!", client.getUsername());
        return new ResponseEntity<>(HttpStatus.OK);
    }

    public Collection<String> findAllUsernames() {
        return clientRepository.findAllUsernames();
    }

    public boolean isPinOk(String username, String pin) {
        Client user = clientRepository.findByUsername(username);
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

    public Client findByUsername(String username) {
        return clientRepository.findByUsername(username);
    }

    public ResponseEntity<?> isPasswordInBlackList(String pass) throws URISyntaxException, IOException {
        Path path = Paths.get(getClass().getClassLoader().getResource("PasswordBlackList.txt").toURI());
        Stream<String> lines = Files.lines(path);
        String data = lines.collect(Collectors.joining("\n"));
        lines.close();
        List<String> passwords = Arrays.asList(data.split("\n"));
        if (passwords.contains(pass))
            return new ResponseEntity<>("Your password has been compromised. Please enter new password.", HttpStatus.OK);
        else
            return new ResponseEntity<>(HttpStatus.OK);

    }

    public ResponseEntity<?> updateClient(Client client, HttpServletRequest request) {
        Authentication authentication = SecurityContextHolder.getContext().getAuthentication();
        Client clientInDb = findByUsername(client.getUsername());
        clientInDb.setEmail(client.getEmail());
        clientInDb.setFirstName(client.getFirstName());
        clientInDb.setLastName(client.getLastName());
        if (!clientInDb.getPassword().equals(client.getPassword()) || client.getPassword() == "") {
            PasswordEncoder passwordEncoder = new BCryptPasswordEncoder();
            clientInDb.setPassword(passwordEncoder.encode(client.getPassword().concat(clientInDb.getSalt())));
            clientInDb.setForgotten(0);
            String pin = RandomStringInitializer.generatePin();
            clientInDb.setPin(passwordEncoder.encode(pin.concat(clientInDb.getSalt())));
        }
        save(clientInDb);
        log.info("Ip: {}, username: {}, User updated successfully!", request.getRemoteAddr(), authentication.getPrincipal().toString());
        return new ResponseEntity<>(clientInDb, HttpStatus.OK);
    }

    public ResponseEntity<?> getByUsername(String username) {
        if (findByUsername(username) == null)
            return new ResponseEntity<>(HttpStatus.BAD_REQUEST);
        return new ResponseEntity<>(findByUsername(username), HttpStatus.OK);
    }

    public ResponseEntity<?> sendPinFor2Auth(String username) {
        return new ResponseEntity<>(HttpStatus.OK);
    }

    public void send2factorAuthPin(Client client) {
        PasswordEncoder passwordEncoder = new BCryptPasswordEncoder();
        String pin = RandomStringInitializer.generatePin();
        client.setPin(passwordEncoder.encode(pin.concat(client.getSalt())));
        client.setPinCreatedDate(new Date());
        clientRepository.save(client);
        emailService.send2factorAuthPin(client.getEmail(), pin);
    }
}
