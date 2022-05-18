package com.example.agent.service;

import com.example.agent.model.*;
import com.example.agent.repository.ClientRepository;
import com.example.agent.repository.CommentRepository;
import com.example.agent.repository.InterviewInformationRepository;
import com.example.agent.repository.JobPositionRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.dao.DataIntegrityViolationException;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;

import java.util.List;
import java.util.Optional;
import java.util.Random;
import java.util.Set;

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

    public ResponseEntity<?> create(Client client) {
        PasswordEncoder passwordEncoder = new BCryptPasswordEncoder();
        try {
            client.setPassword(passwordEncoder.encode(client.getPassword()));
            List<Role> roles = roleService.findByName("ROLE_CLIENT");
            client.setRoles((Set<Role>) roles);
            clientRepository.save(client);
            emailService.sendActivationMailClientAsync(client);
            return new ResponseEntity<>(HttpStatus.OK);
        } catch (DataIntegrityViolationException e) {
            e.printStackTrace();
            return new ResponseEntity<>("Already exist user with same username or email", HttpStatus.BAD_REQUEST);
        }
    }

    public ResponseEntity<?> createComment(Comment comment) {
        commentRepository.save(comment);
        return new ResponseEntity<>(HttpStatus.OK);
    }

    public ResponseEntity<?> updateSalary(JobPosition jobSalary) {
        Optional<JobPosition> job = jobPositionRepository.findById(jobSalary.getId());
        if (!job.isPresent())
            return new ResponseEntity<>(HttpStatus.BAD_REQUEST);
        Random rand = new Random();
        job.get().setAvgSalary(rand.nextInt(2000-1000) + 1000);
        jobPositionRepository.save(job.get());
        return new ResponseEntity<>(HttpStatus.OK);
    }

    public ResponseEntity<?> addInterviewInformation(InterviewInformation interviewInformation) {
        interviewInformationRepository.save(interviewInformation);
        return new ResponseEntity<>(HttpStatus.OK);
    }

    public void save(Client verificationClient) {
        clientRepository.save(verificationClient);
    }
}
