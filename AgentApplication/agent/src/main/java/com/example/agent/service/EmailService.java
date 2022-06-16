package com.example.agent.service;

import com.example.agent.model.Client;
import com.example.agent.model.CompanyOwner;
import com.example.agent.model.ConfirmationToken;
import com.example.agent.repository.ConfirmationTokenRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.core.env.Environment;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.mail.MailException;
import org.springframework.mail.SimpleMailMessage;
import org.springframework.mail.javamail.JavaMailSender;
import org.springframework.scheduling.annotation.Async;
import org.springframework.stereotype.Service;

@Service
public class EmailService {

    @Autowired
    private JavaMailSender javaMailSender;
    @Autowired
    private Environment env;
    @Autowired
    ConfirmationTokenRepository confirmationTokenRepository;


    @Async
    public void sendActivationMailClientAsync(Client user) throws MailException {
        ConfirmationToken confirmationToken = new ConfirmationToken(user);
        confirmationTokenRepository.save(confirmationToken);

        SimpleMailMessage mail = new SimpleMailMessage();
        mail.setTo(user.getEmail());
        mail.setFrom(env.getProperty("spring.mail.username"));
        mail.setSubject("Activation mail");
        mail.setText("Hi, " + user.getFirstName() + ".\n\nWelcome to our site." +
                "\nWe hope that you will be satisfied with our services." +
                "\nIn order to activate your account click on this link: " +
                "https://localhost:8600/clients/activate?token=" + confirmationToken.getConfirmationToken());

        javaMailSender.send(mail);
    }

    @Async
    public void sendActivationMailOwnerAsync(CompanyOwner user) throws MailException {
        ConfirmationToken confirmationToken = new ConfirmationToken(user);
        confirmationTokenRepository.save(confirmationToken);

        SimpleMailMessage mail = new SimpleMailMessage();
        mail.setTo(user.getEmail());
        mail.setFrom(env.getProperty("spring.mail.username"));
        mail.setSubject("Activation mail");
        mail.setText("Hi, " + user.getFirstName() + ".\n\nWelcome to our site." +
                "\nWe hope that you will be satisfied with our services." +
                "\nIn order to activate your account click on this link: " +
                "https://localhost:8600/company/owner/activate?token=" + confirmationToken.getConfirmationToken());

        javaMailSender.send(mail);
    }

    @Async
    public void sendNewPassword(String email, String password) {
        SimpleMailMessage mail = new SimpleMailMessage();
        mail.setTo(email);
        mail.setFrom(env.getProperty("spring.mail.username"));
        mail.setSubject("Refreshed password");
        mail.setText("Your new password is: " + password + ".\nYou have to set your password when you first log in.");

        javaMailSender.send(mail);
    }

    @Async
    public void sendPin(String email, String pin) {
        SimpleMailMessage mail = new SimpleMailMessage();
        mail.setTo(email);
        mail.setFrom(env.getProperty("spring.mail.username"));
        mail.setSubject("New login PIN");
        mail.setText("Your new PIN is: " + pin);

        javaMailSender.send(mail);
    }

    @Async
    public ResponseEntity<?> sendActivationMailToDislinktUser(String email, String name, String key) {
        SimpleMailMessage mail = new SimpleMailMessage();
        mail.setTo(email);
        mail.setFrom(env.getProperty("spring.mail.username"));
        mail.setSubject("Activation mail");
        mail.setText("Hi, " + name + ".\n\nWelcome to our site." +
                "\nWe hope that you will be satisfied with our services." +
                "\nIn order to activate your account click on this link: " +
                "https://localhost:8000/user/activate?token=" + key);

        javaMailSender.send(mail);
        return new ResponseEntity<>(HttpStatus.OK);
    }

    @Async
    public ResponseEntity<?> sendPasswordless(String email, String salt) {
        SimpleMailMessage mail = new SimpleMailMessage();
        mail.setTo(email);
        mail.setFrom(env.getProperty("spring.mail.username"));
        mail.setSubject("Passwordless authentication");
        mail.setText("In order to login click on this link: " +
                "https://localhost:4200/passwordless?token=" + salt);

        javaMailSender.send(mail);
        return new ResponseEntity<>(HttpStatus.OK);
    }

    @Async
    public ResponseEntity<?> sendPasswordlesstoDislinkt(String email, String salt) {
        SimpleMailMessage mail = new SimpleMailMessage();
        mail.setTo(email);
        mail.setFrom(env.getProperty("spring.mail.username"));
        mail.setSubject("Passwordless authentication");
        mail.setText("In order to login click on this link: " +
                "https://localhost:4300/passwordless?token=" + salt);

        javaMailSender.send(mail);
        return new ResponseEntity<>(HttpStatus.OK);
    }

    @Async
    public void send2factorAuthPin(String email, String pin) {
        SimpleMailMessage mail = new SimpleMailMessage();
        mail.setTo(email);
        mail.setFrom(env.getProperty("spring.mail.username"));
        mail.setSubject("2 factor authentication PIN");
        mail.setText("Your PIN is: " + pin);

        javaMailSender.send(mail);
    }
}
