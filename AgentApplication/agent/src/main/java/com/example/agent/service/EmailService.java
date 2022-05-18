package com.example.agent.service;

import com.example.agent.model.Client;
import com.example.agent.model.CompanyOwner;
import com.example.agent.model.ConfirmationToken;
import com.example.agent.model.User;
import com.example.agent.repository.ConfirmationTokenRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.core.env.Environment;
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
                "http://localhost:8600/clients/activate?token=" + confirmationToken.getConfirmationToken());

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
                "http://localhost:8600/clients/activate?token=" + confirmationToken.getConfirmationToken());

        javaMailSender.send(mail);
    }

}
