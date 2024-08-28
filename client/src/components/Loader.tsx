import { Spinner, Center, Box } from '@chakra-ui/react';

const Loader = () => {
    return (
        <Box
            position="fixed"
            top="0"
            left="0"
            width="100%"
            height="100%"
            bg="rgba(255, 255, 255, 0.8)"
            backdropFilter="blur(5px)"
            zIndex="1000"
        >
            <Center height="100%">
                <Spinner size="xl" />
            </Center>
        </Box>
    );
}

export default Loader;