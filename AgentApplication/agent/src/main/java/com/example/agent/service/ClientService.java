package com.example.agent.service;

import com.example.agent.model.*;
import com.example.agent.model.dto.UserDto;
import com.example.agent.repository.*;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.dao.DataIntegrityViolationException;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;

import java.io.IOException;
import java.net.URISyntaxException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.*;
import java.util.stream.Collectors;
import java.util.stream.Stream;

@Service
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

    public ResponseEntity<?> create(UserDto dto) {
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
            return new ResponseEntity<>(HttpStatus.OK);
        } catch (DataIntegrityViolationException e) {
            e.printStackTrace();
            return new ResponseEntity<>("Already exist user with same username or email", HttpStatus.BAD_REQUEST);
        }
    }

    public ResponseEntity<?> createComment(Comment comment, Integer companyId) {
        Company company = companyRepository.getById(companyId);
        comment.setCompany(company);
        commentRepository.save(comment);
        return new ResponseEntity<>(HttpStatus.OK);
    }

    public ResponseEntity<?> updateSalary(JobPosition jobSalary) {
        Optional<JobPosition> job = jobPositionRepository.findById(jobSalary.getId());
        if (!job.isPresent())
            return new ResponseEntity<>(HttpStatus.BAD_REQUEST);
        Random rand = new Random();
        job.get().setAvgSalary(rand.nextInt(2000 - 1000) + 1000);
        jobPositionRepository.save(job.get());
        return new ResponseEntity<>(HttpStatus.OK);
    }

    public ResponseEntity<?> addInterviewInformation(InterviewInformation interviewInformation, Integer jobId) {
        interviewInformation.setJobPosition(jobPositionRepository.getById(jobId));
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
        return new ResponseEntity<>(HttpStatus.OK);
    }

    public Collection<String> findAllUsernames() {
        return clientRepository.findAllUsernames();
    }

    public boolean isPinOk(String username, String pin) {
        Client user = clientRepository.findByUsername(username);
        if (user == null)
            return false;
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

    public ResponseEntity<?> updateClient(Client client) {
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
        return new ResponseEntity<>(clientInDb, HttpStatus.OK);
    }

    public ResponseEntity<?> getByUsername(String username) {
        if (findByUsername(username) == null)
            return new ResponseEntity<>(HttpStatus.BAD_REQUEST);
        return new ResponseEntity<>(findByUsername(username), HttpStatus.OK);
    }
}
