import React, { useState } from 'react';
import { Box, Button, FormControl, FormLabel, Input, Heading } from '@chakra-ui/react';

const Register: React.FC = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    // Aqui você pode adicionar a lógica de registro
    console.log('Registro:', { email, password, confirmPassword });
  };

  return (
    <Box
        p={4} 
        borderWidth={1} 
        borderRadius="lg" 
        display="flex" 
        flexDirection="column" 
        justifyContent="center" 
        alignItems="center"
    >
      <Heading as="h2" size="lg" mb={4}>Registro</Heading>
      <form onSubmit={handleSubmit}>
        <FormControl mb={4} isRequired>
          <FormLabel>Email</FormLabel>
          <Input
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            placeholder="Digite seu email"
            w={300}
          />
        </FormControl>
        <FormControl mb={4} isRequired>
          <FormLabel>Senha</FormLabel>
          <Input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            placeholder="Digite sua senha"
            w={300}
          />
        </FormControl>
        <FormControl mb={4} isRequired>
          <FormLabel>Confirmar Senha</FormLabel>
          <Input
            type="password"
            value={confirmPassword}
            onChange={(e) => setConfirmPassword(e.target.value)}
            placeholder="Confirme sua senha"
            w={300}
          />
        </FormControl>
        <Button type="submit" colorScheme="teal" width="full">Registrar</Button>
      </form>
    </Box>
  );
};

export default Register;