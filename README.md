# Barebone-Container-Implemenation-in-Go

## Implement a barebones container.

POC idea: 
Build a basic container runtime with chroot, implementing isolation in networking, security, process management.
Build it all from scratch, in golang. 
Key goal: It should be able to run a service inside an isolated "environment" via tricking the user. Chroot with steroids.

## Tests

```bash
# In the project root file
container_test.go
```
### Explanation:
    TestRunCommand: Tests the run command by executing a simple echo command.
    TestChildCommand: Verifies the child command executes and produces the correct output.
    TestSetupCgroups: Creates a temporary directory to mock cgroup paths and verifies the files are created with the expected values. It uses Go's t.TempDir for an isolated environment.

### Notes:
    Ensure to run the tests with sufficient/relevant permissions, as some system-level operations may require root.
    Mocking real cgroup paths with a temporary directory avoids affecting the host system.
    This assumes /proc/self/exe works as expected in the test environment. Keep this in mind!

#### Inspiriations and sources: 

- https://www.youtube.com/watch?v=Utf-A4rODH8  Liz Rice - Building a container from scratch in Go.
- https://www.youtube.com/watch?v=JOsWB50LmwQ Earthly - Build your own Container Runtime.
- https://www.youtube.com/watch?v=sK5i-N34im8 Jérôme Petazzoni - Cgroups, namespaces, and beyond: what are containers made from? 
- https://www.youtube.com/watch?v=GFUpXhft8zA Vishesh @DeepSource - Building containers from scratch 
